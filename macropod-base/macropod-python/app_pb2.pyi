from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class RequestBody(_message.Message):
    __slots__ = ("data", "workflow_id", "depth", "width", "request_type", "pv_path")
    DATA_FIELD_NUMBER: _ClassVar[int]
    WORKFLOW_ID_FIELD_NUMBER: _ClassVar[int]
    DEPTH_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    REQUEST_TYPE_FIELD_NUMBER: _ClassVar[int]
    PV_PATH_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    workflow_id: str
    depth: int
    width: int
    request_type: str
    pv_path: str
    def __init__(self, data: _Optional[bytes] = ..., workflow_id: _Optional[str] = ..., depth: _Optional[int] = ..., width: _Optional[int] = ..., request_type: _Optional[str] = ..., pv_path: _Optional[str] = ...) -> None: ...

class ResponseBody(_message.Message):
    __slots__ = ("reply", "code", "pv_path")
    REPLY_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    PV_PATH_FIELD_NUMBER: _ClassVar[int]
    reply: str
    code: int
    pv_path: str
    def __init__(self, reply: _Optional[str] = ..., code: _Optional[int] = ..., pv_path: _Optional[str] = ...) -> None: ...
