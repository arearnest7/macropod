package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"flag"
	"fmt"
	"net"
	"os"
	"io/ioutil"
	"log"
	"os/exec"

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



// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/video-decoder")
	if err != nil {
      		log.Fatal(err)
   	}
	video-decoder := string(contents)

	videoFile = "reference/video.mp4"
	videoFragment, _ = os.ReadFile(*videoFile)
	log.Infof("read video fragment, size: %v\n", len(videoFragment))

	requestURL := video-decoder + ":80"
	req_url, err := http.NewRequest(http.MethodPost, requestURL, {"video": videoFragment})
	if err != nil {
		log.Fatal(err)
	}
	ret, err := exec.Command("python3 decoder.py " + req_url).Output()
	fmt.Fprintf(res, ret) // echo to caller
}
