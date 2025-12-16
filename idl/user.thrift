namespace go user

include "base.thrift"

struct User {
  1: i64 id
  2: string name
  3: string email
  4: string yob // Year of birth
}

struct GetUserReq {
  1: i64 id
}

struct GetUserResp {
  1: User user

  255: base.BaseResp baseResp
}

service UserService {
  GetUserResp GetUser(1: GetUserReq req)
}