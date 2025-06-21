from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MacroPodRequest(_message.Message):
    __slots__ = ("Text", "JSON", "Data", "Workflow", "Function", "WorkflowID", "Depth", "Width", "Target")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    JSON_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    WORKFLOW_FIELD_NUMBER: _ClassVar[int]
    FUNCTION_FIELD_NUMBER: _ClassVar[int]
    WORKFLOWID_FIELD_NUMBER: _ClassVar[int]
    DEPTH_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    TARGET_FIELD_NUMBER: _ClassVar[int]
    Text: str
    JSON: _struct_pb2.Struct
    Data: bytes
    Workflow: str
    Function: str
    WorkflowID: str
    Depth: int
    Width: int
    Target: str
    def __init__(self, Text: _Optional[str] = ..., JSON: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., Data: _Optional[bytes] = ..., Workflow: _Optional[str] = ..., Function: _Optional[str] = ..., WorkflowID: _Optional[str] = ..., Depth: _Optional[int] = ..., Width: _Optional[int] = ..., Target: _Optional[str] = ...) -> None: ...

class MacroPodReply(_message.Message):
    __slots__ = ("Reply", "Code")
    REPLY_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    Reply: str
    Code: int
    def __init__(self, Reply: _Optional[str] = ..., Code: _Optional[int] = ...) -> None: ...

class FunctionStruct(_message.Message):
    __slots__ = ("Registry", "Endpoints", "Envs", "Secrets")
    class EndpointsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class EnvsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class SecretsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    REGISTRY_FIELD_NUMBER: _ClassVar[int]
    ENDPOINTS_FIELD_NUMBER: _ClassVar[int]
    ENVS_FIELD_NUMBER: _ClassVar[int]
    SECRETS_FIELD_NUMBER: _ClassVar[int]
    Registry: str
    Endpoints: _containers.ScalarMap[str, str]
    Envs: _containers.ScalarMap[str, str]
    Secrets: _containers.ScalarMap[str, str]
    def __init__(self, Registry: _Optional[str] = ..., Endpoints: _Optional[_Mapping[str, str]] = ..., Envs: _Optional[_Mapping[str, str]] = ..., Secrets: _Optional[_Mapping[str, str]] = ...) -> None: ...

class ConfigStruct(_message.Message):
    __slots__ = ("Namespace", "TTL", "Deployment", "Communication", "Aggregation", "TargetConcurrency", "Debug")
    NAMESPACE_FIELD_NUMBER: _ClassVar[int]
    TTL_FIELD_NUMBER: _ClassVar[int]
    DEPLOYMENT_FIELD_NUMBER: _ClassVar[int]
    COMMUNICATION_FIELD_NUMBER: _ClassVar[int]
    AGGREGATION_FIELD_NUMBER: _ClassVar[int]
    TARGETCONCURRENCY_FIELD_NUMBER: _ClassVar[int]
    DEBUG_FIELD_NUMBER: _ClassVar[int]
    Namespace: str
    TTL: int
    Deployment: str
    Communication: str
    Aggregation: str
    TargetConcurrency: int
    Debug: int
    def __init__(self, Namespace: _Optional[str] = ..., TTL: _Optional[int] = ..., Deployment: _Optional[str] = ..., Communication: _Optional[str] = ..., Aggregation: _Optional[str] = ..., TargetConcurrency: _Optional[int] = ..., Debug: _Optional[int] = ...) -> None: ...

class PayloadStruct(_message.Message):
    __slots__ = ("Type", "Text", "JSON", "Data")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    JSON_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    Type: str
    Text: str
    JSON: _struct_pb2.Struct
    Data: bytes
    def __init__(self, Type: _Optional[str] = ..., Text: _Optional[str] = ..., JSON: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., Data: _Optional[bytes] = ...) -> None: ...

class WorkflowStruct(_message.Message):
    __slots__ = ("Name", "Functions", "Config", "Payload")
    class FunctionsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: FunctionStruct
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[FunctionStruct, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    FUNCTIONS_FIELD_NUMBER: _ClassVar[int]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_FIELD_NUMBER: _ClassVar[int]
    Name: str
    Functions: _containers.MessageMap[str, FunctionStruct]
    Config: ConfigStruct
    Payload: PayloadStruct
    def __init__(self, Name: _Optional[str] = ..., Functions: _Optional[_Mapping[str, FunctionStruct]] = ..., Config: _Optional[_Union[ConfigStruct, _Mapping]] = ..., Payload: _Optional[_Union[PayloadStruct, _Mapping]] = ...) -> None: ...

class EvalStruct(_message.Message):
    __slots__ = ("WorkflowConcurrency", "Invocations", "ExtraTargets", "ExtraTargetsPayload", "Workflows")
    class ExtraTargetsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class ExtraTargetsPayloadEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: PayloadStruct
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[PayloadStruct, _Mapping]] = ...) -> None: ...
    class WorkflowsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: WorkflowStruct
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[WorkflowStruct, _Mapping]] = ...) -> None: ...
    WORKFLOWCONCURRENCY_FIELD_NUMBER: _ClassVar[int]
    INVOCATIONS_FIELD_NUMBER: _ClassVar[int]
    EXTRATARGETS_FIELD_NUMBER: _ClassVar[int]
    EXTRATARGETSPAYLOAD_FIELD_NUMBER: _ClassVar[int]
    WORKFLOWS_FIELD_NUMBER: _ClassVar[int]
    WorkflowConcurrency: _containers.RepeatedScalarFieldContainer[int]
    Invocations: int
    ExtraTargets: _containers.ScalarMap[str, str]
    ExtraTargetsPayload: _containers.MessageMap[str, PayloadStruct]
    Workflows: _containers.MessageMap[str, WorkflowStruct]
    def __init__(self, WorkflowConcurrency: _Optional[_Iterable[int]] = ..., Invocations: _Optional[int] = ..., ExtraTargets: _Optional[_Mapping[str, str]] = ..., ExtraTargetsPayload: _Optional[_Mapping[str, PayloadStruct]] = ..., Workflows: _Optional[_Mapping[str, WorkflowStruct]] = ...) -> None: ...

class MetricsStruct(_message.Message):
    __slots__ = ("Uptime", "LoadAvg1", "LoadAvg5", "LoadAvg15", "CPUUsed", "CPUCountLogical", "CPUCountPhysical", "MemoryUsed", "MemoryAvailable", "MemoryTotal", "MemoryBuffers", "MemoryCached", "MemoryWriteBack", "MemoryDirty", "MemoryWriteBackTmp", "MemoryShared", "MemorySlab", "MemorySreclaimable", "MemorySunreclaim", "MemoryPageTables", "MemorySwapCached", "MemoryCommitLimit", "MemoryCommittedAS", "MemoryHighTotal", "MemoryHighFree", "MemoryLowTotal", "MemoryLowFree", "MemorySwapTotal", "MemorySwapFree", "MemoryMapped", "MemoryVmallocTotal", "MemoryVmallocUsed", "MemoryVmallocChunk", "MemoryHugePagesTotal", "MemoryHugePagesFree", "MemoryHugePagesRsvd", "MemoryHugePagesSurp", "MemoryHugePageSize", "MemoryAnonHugePages", "DiskUsed", "DiskFree", "DiskTotal", "DiskInodesUsed", "DiskInodesFree", "DiskInodesTotal", "NetworkBytesSent", "NetworkBytesRecv", "NetworkPacketsSent", "NetworkPacketsRecv", "NetworkErrin", "NetworkErrout", "NetworkDropin", "NetworkDropout", "NetworkFifoin", "NetworkFifoout")
    UPTIME_FIELD_NUMBER: _ClassVar[int]
    LOADAVG1_FIELD_NUMBER: _ClassVar[int]
    LOADAVG5_FIELD_NUMBER: _ClassVar[int]
    LOADAVG15_FIELD_NUMBER: _ClassVar[int]
    CPUUSED_FIELD_NUMBER: _ClassVar[int]
    CPUCOUNTLOGICAL_FIELD_NUMBER: _ClassVar[int]
    CPUCOUNTPHYSICAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYUSED_FIELD_NUMBER: _ClassVar[int]
    MEMORYAVAILABLE_FIELD_NUMBER: _ClassVar[int]
    MEMORYTOTAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYBUFFERS_FIELD_NUMBER: _ClassVar[int]
    MEMORYCACHED_FIELD_NUMBER: _ClassVar[int]
    MEMORYWRITEBACK_FIELD_NUMBER: _ClassVar[int]
    MEMORYDIRTY_FIELD_NUMBER: _ClassVar[int]
    MEMORYWRITEBACKTMP_FIELD_NUMBER: _ClassVar[int]
    MEMORYSHARED_FIELD_NUMBER: _ClassVar[int]
    MEMORYSLAB_FIELD_NUMBER: _ClassVar[int]
    MEMORYSRECLAIMABLE_FIELD_NUMBER: _ClassVar[int]
    MEMORYSUNRECLAIM_FIELD_NUMBER: _ClassVar[int]
    MEMORYPAGETABLES_FIELD_NUMBER: _ClassVar[int]
    MEMORYSWAPCACHED_FIELD_NUMBER: _ClassVar[int]
    MEMORYCOMMITLIMIT_FIELD_NUMBER: _ClassVar[int]
    MEMORYCOMMITTEDAS_FIELD_NUMBER: _ClassVar[int]
    MEMORYHIGHTOTAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYHIGHFREE_FIELD_NUMBER: _ClassVar[int]
    MEMORYLOWTOTAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYLOWFREE_FIELD_NUMBER: _ClassVar[int]
    MEMORYSWAPTOTAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYSWAPFREE_FIELD_NUMBER: _ClassVar[int]
    MEMORYMAPPED_FIELD_NUMBER: _ClassVar[int]
    MEMORYVMALLOCTOTAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYVMALLOCUSED_FIELD_NUMBER: _ClassVar[int]
    MEMORYVMALLOCCHUNK_FIELD_NUMBER: _ClassVar[int]
    MEMORYHUGEPAGESTOTAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYHUGEPAGESFREE_FIELD_NUMBER: _ClassVar[int]
    MEMORYHUGEPAGESRSVD_FIELD_NUMBER: _ClassVar[int]
    MEMORYHUGEPAGESSURP_FIELD_NUMBER: _ClassVar[int]
    MEMORYHUGEPAGESIZE_FIELD_NUMBER: _ClassVar[int]
    MEMORYANONHUGEPAGES_FIELD_NUMBER: _ClassVar[int]
    DISKUSED_FIELD_NUMBER: _ClassVar[int]
    DISKFREE_FIELD_NUMBER: _ClassVar[int]
    DISKTOTAL_FIELD_NUMBER: _ClassVar[int]
    DISKINODESUSED_FIELD_NUMBER: _ClassVar[int]
    DISKINODESFREE_FIELD_NUMBER: _ClassVar[int]
    DISKINODESTOTAL_FIELD_NUMBER: _ClassVar[int]
    NETWORKBYTESSENT_FIELD_NUMBER: _ClassVar[int]
    NETWORKBYTESRECV_FIELD_NUMBER: _ClassVar[int]
    NETWORKPACKETSSENT_FIELD_NUMBER: _ClassVar[int]
    NETWORKPACKETSRECV_FIELD_NUMBER: _ClassVar[int]
    NETWORKERRIN_FIELD_NUMBER: _ClassVar[int]
    NETWORKERROUT_FIELD_NUMBER: _ClassVar[int]
    NETWORKDROPIN_FIELD_NUMBER: _ClassVar[int]
    NETWORKDROPOUT_FIELD_NUMBER: _ClassVar[int]
    NETWORKFIFOIN_FIELD_NUMBER: _ClassVar[int]
    NETWORKFIFOOUT_FIELD_NUMBER: _ClassVar[int]
    Uptime: float
    LoadAvg1: float
    LoadAvg5: float
    LoadAvg15: float
    CPUUsed: float
    CPUCountLogical: float
    CPUCountPhysical: float
    MemoryUsed: float
    MemoryAvailable: float
    MemoryTotal: float
    MemoryBuffers: float
    MemoryCached: float
    MemoryWriteBack: float
    MemoryDirty: float
    MemoryWriteBackTmp: float
    MemoryShared: float
    MemorySlab: float
    MemorySreclaimable: float
    MemorySunreclaim: float
    MemoryPageTables: float
    MemorySwapCached: float
    MemoryCommitLimit: float
    MemoryCommittedAS: float
    MemoryHighTotal: float
    MemoryHighFree: float
    MemoryLowTotal: float
    MemoryLowFree: float
    MemorySwapTotal: float
    MemorySwapFree: float
    MemoryMapped: float
    MemoryVmallocTotal: float
    MemoryVmallocUsed: float
    MemoryVmallocChunk: float
    MemoryHugePagesTotal: float
    MemoryHugePagesFree: float
    MemoryHugePagesRsvd: float
    MemoryHugePagesSurp: float
    MemoryHugePageSize: float
    MemoryAnonHugePages: float
    DiskUsed: float
    DiskFree: float
    DiskTotal: float
    DiskInodesUsed: float
    DiskInodesFree: float
    DiskInodesTotal: float
    NetworkBytesSent: float
    NetworkBytesRecv: float
    NetworkPacketsSent: float
    NetworkPacketsRecv: float
    NetworkErrin: float
    NetworkErrout: float
    NetworkDropin: float
    NetworkDropout: float
    NetworkFifoin: float
    NetworkFifoout: float
    def __init__(self, Uptime: _Optional[float] = ..., LoadAvg1: _Optional[float] = ..., LoadAvg5: _Optional[float] = ..., LoadAvg15: _Optional[float] = ..., CPUUsed: _Optional[float] = ..., CPUCountLogical: _Optional[float] = ..., CPUCountPhysical: _Optional[float] = ..., MemoryUsed: _Optional[float] = ..., MemoryAvailable: _Optional[float] = ..., MemoryTotal: _Optional[float] = ..., MemoryBuffers: _Optional[float] = ..., MemoryCached: _Optional[float] = ..., MemoryWriteBack: _Optional[float] = ..., MemoryDirty: _Optional[float] = ..., MemoryWriteBackTmp: _Optional[float] = ..., MemoryShared: _Optional[float] = ..., MemorySlab: _Optional[float] = ..., MemorySreclaimable: _Optional[float] = ..., MemorySunreclaim: _Optional[float] = ..., MemoryPageTables: _Optional[float] = ..., MemorySwapCached: _Optional[float] = ..., MemoryCommitLimit: _Optional[float] = ..., MemoryCommittedAS: _Optional[float] = ..., MemoryHighTotal: _Optional[float] = ..., MemoryHighFree: _Optional[float] = ..., MemoryLowTotal: _Optional[float] = ..., MemoryLowFree: _Optional[float] = ..., MemorySwapTotal: _Optional[float] = ..., MemorySwapFree: _Optional[float] = ..., MemoryMapped: _Optional[float] = ..., MemoryVmallocTotal: _Optional[float] = ..., MemoryVmallocUsed: _Optional[float] = ..., MemoryVmallocChunk: _Optional[float] = ..., MemoryHugePagesTotal: _Optional[float] = ..., MemoryHugePagesFree: _Optional[float] = ..., MemoryHugePagesRsvd: _Optional[float] = ..., MemoryHugePagesSurp: _Optional[float] = ..., MemoryHugePageSize: _Optional[float] = ..., MemoryAnonHugePages: _Optional[float] = ..., DiskUsed: _Optional[float] = ..., DiskFree: _Optional[float] = ..., DiskTotal: _Optional[float] = ..., DiskInodesUsed: _Optional[float] = ..., DiskInodesFree: _Optional[float] = ..., DiskInodesTotal: _Optional[float] = ..., NetworkBytesSent: _Optional[float] = ..., NetworkBytesRecv: _Optional[float] = ..., NetworkPacketsSent: _Optional[float] = ..., NetworkPacketsRecv: _Optional[float] = ..., NetworkErrin: _Optional[float] = ..., NetworkErrout: _Optional[float] = ..., NetworkDropin: _Optional[float] = ..., NetworkDropout: _Optional[float] = ..., NetworkFifoin: _Optional[float] = ..., NetworkFifoout: _Optional[float] = ...) -> None: ...
