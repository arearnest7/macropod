// MIT License
//
// Copyright (c) 2021 Michal Baczun and EASE lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc/credentials/insecure"

	ctrdlog "github.com/containerd/containerd/log"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb_video "tests/video_analytics/proto"

	sdk "github.com/ease-lab/vhive-xdt/sdk/golang"
	"github.com/ease-lab/vhive-xdt/utils"
	pb_helloworld "github.com/vhive-serverless/vSwarm/examples/protobuf/helloworld"

	storage "github.com/vhive-serverless/vSwarm/utils/storage/go"
	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

var (
	videoFragment []byte
	videoFile     *string
	AWS_S3_BUCKET = "vhive-video-bench"
)

const (
	INLINE = "INLINE"
	XDT    = "XDT"
	S3     = "S3"
)

type server struct {
	decoderAddr    string
	decoderPort    int
	transferType   string
	config         utils.Config
	XDTclient      *sdk.XDTclient
	storageBackend storage.Storage
	pb_helloworld.UnimplementedGreeterServer
}

func fetchSelfIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("Error fetching self IP: " + err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	log.Errorf("unable to find IP, returning empty string")
	return ""
}

func uploadToS3(ctx context.Context, storageBackend storage.Storage) {
	if tracing.IsTracingEnabled() {
		span := tracing.Span{SpanName: "Video upload", TracerName: "S3 video upload - tracer"}
		span.StartSpan(ctx)
		defer span.EndSpan()
	}
	file, err := os.Open(*videoFile)
	if err != nil {
		log.Fatalf("[Video Streaming] Failed to open file: %s", err)
	}
	storageBackend.PutFile("streaming-video.mp4", file)
	log.Infof("[Video Streaming] Uploaded video to s3")
}

// SayHello implements the helloworld interface. Used to trigger the video streamer to start the benchmark.
func (s *server) SayHello(ctx context.Context, req *pb_helloworld.HelloRequest) (_ *pb_helloworld.HelloReply, err error) {
	// Become a client of the decoder function and send the video:
	// establish a connection
	addr := fmt.Sprintf("%v:%v", s.decoderAddr, s.decoderPort)
	log.Infof("[Video Streaming] Using addr: %v", addr)

	// send message
	log.Infof("[Video Streaming] Video Fragment length: %v", len(videoFragment))

	var reply *pb_video.DecodeReply
	var response string
	if s.transferType == XDT {
		payloadToSend := utils.Payload{
			FunctionName: "HelloXDT",
			Data:         videoFragment,
		}
		if message, _, err := (*s.XDTclient).Invoke(ctx, addr, payloadToSend); err != nil {
			log.Fatalf("SQP_to_dQP_data_transfer failed %v", err)
		} else {
			response = string(message)
		}
	} else if s.transferType == S3 || s.transferType == INLINE {
		var conn *grpc.ClientConn
		if tracing.IsTracingEnabled() {
			conn, err = tracing.DialGRPCWithUnaryInterceptor(addr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			conn, err = grpc.Dial(addr, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
		if err != nil {
			log.Fatalf("[Video Streaming] Failed to dial decoder: %s", err)
		}
		defer conn.Close()

		client := pb_video.NewVideoDecoderClient(conn)
		if s.transferType == S3 {
			// upload video to s3
			uploadToS3(ctx, s.storageBackend)
			// issue request
			reply, err = client.Decode(ctx, &pb_video.DecodeRequest{S3Key: "streaming-video.mp4"})
		} else {
			reply, err = client.Decode(ctx, &pb_video.DecodeRequest{Video: videoFragment})
		}
		if err != nil {
			log.Fatalf("[Video Streaming] Failed to send video to decoder: %s", err)
		}
		response = reply.Classification
	} else {
		log.Fatalf("Invalid TRANSFER_TYPE value")
	}

	log.Infof("[Video Streaming] Received Decoder reply")
	return &pb_helloworld.HelloReply{Message: response}, err
}

func main() {
	debug := flag.Bool("d", false, "Debug level in logs")
	dockerCompose := flag.Bool("dockerCompose", false, "Execution env")
	decoderAddr := flag.String("addr", "decoder.default.svc.cluster.local", "Decoder address")
	decoderPort := flag.Int("p", 80, "Decoder port")
	servePort := flag.Int("sp", 80, "Port listened to by this streamer")
	videoFile = flag.String("video", "reference/video.mp4", "The file location of the video")
	zipkin := flag.String("zipkin", "http://zipkin.istio-system.svc.cluster.local:9411/api/v2/spans", "zipkin url")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})

	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging is enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
	if tracing.IsTracingEnabled() {
		shutdown, err := tracing.InitBasicTracer(*zipkin, "Video Streaming")
		if err != nil {
			log.Warn(err)
		}
		defer shutdown()
	}
	videoFragment, _ = os.ReadFile(*videoFile)
	log.Infof("read video fragment, size: %v\n", len(videoFragment))
	// server setup: listen on port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *servePort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcServer *grpc.Server
	if tracing.IsTracingEnabled() {
		grpcServer = tracing.GetGRPCServerWithUnaryInterceptor()
	} else {
		grpcServer = grpc.NewServer()
	}

	reflection.Register(grpcServer)
	server := server{}
	server.decoderAddr = *decoderAddr
	server.decoderPort = *decoderPort

	server.transferType = INLINE
	if transferType, ok := os.LookupEnv("TRANSFER_TYPE"); !ok {
		server.transferType = INLINE
	} else {
		server.transferType = transferType
	}

	if server.transferType == S3 {
		if value, ok := os.LookupEnv("BUCKET_NAME"); ok {
			AWS_S3_BUCKET = value
		}
		log.Infof("[streaming]  BUCKET = %s", AWS_S3_BUCKET)
		storageBackend := storage.New("S3", AWS_S3_BUCKET)
		server.storageBackend = storageBackend
	} else if server.transferType == XDT {
		log.Infof("[streaming] TransferType = %s", server.transferType)
		config := utils.ReadConfig()
		log.Info(config)
		if !*dockerCompose {
			config.SQPServerHostname = fetchSelfIP()
		}
		xdtClient, err := sdk.NewXDTclient(config)
		if err != nil {
			log.Fatalf("InitXDT failed %v", err)
		}

		server.config = config
		server.XDTclient = xdtClient
		log.Infof("[streaming] XDT client created")
	}
	pb_helloworld.RegisterGreeterServer(grpcServer, &server)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
