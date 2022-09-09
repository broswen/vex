import {getToken, handleRequest} from "@/index";


test("should get bearer token", () => {
  const request = new Request("https://test.com")
  request.headers.set('Authorization', 'Bearer abc123')
  const token = getToken(request)
  expect(token).toEqual('abc123')
})

test("should not get bearer token", () => {
  const request = new Request("https://test.com")
  const token = getToken(request)
  expect(token).toBeNull()
})