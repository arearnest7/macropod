from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

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
    Endpoints: _containers.RepeatedScalarFieldContainer[str]
    Envs: _containers.ScalarMap[str, str]
    Secrets: _containers.ScalarMap[str, str]
    def __init__(self, Registry: _Optional[str] = ..., Endpoints: _Optional[_Iterable[str]] = ..., Envs: _Optional[_Mapping[str, str]] = ..., Secrets: _Optional[_Mapping[str, str]] = ...) -> None: ...

class ConfigStruct(_message.Message):
    __slots__ = ("Namespace", "TTL", "Deployment", "Communication", "Aggregation", "Channel", "TargetConcurrency", "Debug")
    NAMESPACE_FIELD_NUMBER: _ClassVar[int]
    TTL_FIELD_NUMBER: _ClassVar[int]
    DEPLOYMENT_FIELD_NUMBER: _ClassVar[int]
    COMMUNICATION_FIELD_NUMBER: _ClassVar[int]
    AGGREGATION_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_FIELD_NUMBER: _ClassVar[int]
    TARGETCONCURRENCY_FIELD_NUMBER: _ClassVar[int]
    DEBUG_FIELD_NUMBER: _ClassVar[int]
    Namespace: str
    TTL: int
    Deployment: str
    Communication: str
    Aggregation: str
    Channel: str
    TargetConcurrency: int
    Debug: int
    def __init__(self, Namespace: _Optional[str] = ..., TTL: _Optional[int] = ..., Deployment: _Optional[str] = ..., Communication: _Optional[str] = ..., Aggregation: _Optional[str] = ..., Channel: _Optional[str] = ..., TargetConcurrency: _Optional[int] = ..., Debug: _Optional[int] = ...) -> None: ...

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
    __slots__ = ("CPUUsed", "CPUCountPhysical", "MemoryUsed", "NetworkBytesSent")
    CPUUSED_FIELD_NUMBER: _ClassVar[int]
    CPUCOUNTPHYSICAL_FIELD_NUMBER: _ClassVar[int]
    MEMORYUSED_FIELD_NUMBER: _ClassVar[int]
    NETWORKBYTESSENT_FIELD_NUMBER: _ClassVar[int]
    CPUUsed: float
    CPUCountPhysical: float
    MemoryUsed: float
    NetworkBytesSent: float
    def __init__(self, CPUUsed: _Optional[float] = ..., CPUCountPhysical: _Optional[float] = ..., MemoryUsed: _Optional[float] = ..., NetworkBytesSent: _Optional[float] = ...) -> None: ...
