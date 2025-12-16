namespace go yob

include "base.thrift"

// Year of birth
struct GetUserYOBReq {
  1: i64 user_id
}

struct GetUserYOBResp {
  1: i64 yob

  255: base.BaseResp baseResp
}

service YOBService {
  GetUserYOBResp GetUserYOB(1: GetUserYOBReq req)
}

