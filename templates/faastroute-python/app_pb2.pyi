from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class RequestBody(_message.Message):
    __slots__ = ("body", "workflow_id")
    BODY_FIELD_NUMBER: _ClassVar[int]
    WORKFLOW_ID_FIELD_NUMBER: _ClassVar[int]
    body: bytes
    workflow_id: str
    def __init__(self, body: _Optional[bytes] = ..., workflow_id: _Optional[str] = ...) -> None: ...

class ResponseBody(_message.Message):
    __slots__ = ("reply", "code")
    REPLY_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    reply: str
    code: int
    def __init__(self, reply: _Optional[str] = ..., code: _Optional[int] = ...) -> None: ...
