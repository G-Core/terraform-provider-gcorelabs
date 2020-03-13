provider "gcore" {
  jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoyNTgwOTM4Nzg4LCJqdGkiOiJhZWQ2ZjQwNjhhYzM0NWNkYWM1MTcwZjk0MzcwMDIzMyIsInVzZXJfaWQiOjEsInVzZXJfdHlwZSI6InN5c3RlbV9hZG1pbiIsInVzZXJfZ3JvdXBzIjpudWxsLCJjbGllbnRfaWQiOm51bGwsImVtYWlsIjoidGVzdEB0ZXN0LnRlc3QiLCJ1c2VybmFtZSI6InRlc3RAdGVzdC50ZXN0IiwiaXNfYWRtaW4iOnRydWUsImNsaWVudF9uYW1lIjoidGVzdCIsInJqdGkiOiJjZmJhNzMxODhlOTg0MzgxODAzZDdmYzU3OWJmZWIxYyJ9.0xjny_NM1uLQ5gRT8ZSmA_tvyeNZs8BrPjSFfhkKJbk"
}

resource "gcore_volume" "foo" {
  name = 156
  size = 2
  type_name = "ssd_hiiops"
  region_id = 1
  project_id = 78
  count = 2
}