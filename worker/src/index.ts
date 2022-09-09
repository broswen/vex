
export type Env = {
  FLAG: KVNamespace
  TOKEN: KVNamespace
}

export async function handleRequest(request: Request, env: Env) {
  const url = new URL(request.url);
  const projectId = url.pathname.slice(1)
  const token = getToken(request)

  //reject if no bearer token
  if (!token) {
    return new Response('missing token', {status: 401})
  }

  //reject if token not in kv
  const tokenAccount = await getTokenAccount(token, env)
  if (!tokenAccount) {
    return new Response('invalid token', {status: 401})
  }

  //reject if project id is invalid
  if (projectId.length !== 36) {
    return new Response('invalid project id', {status: 400})
  }
  const getWithMetadataResult = await env.FLAG.getWithMetadata(projectId)

  //reject if bearer token account id doesn't match project account id from metadata
  if (getWithMetadataResult.metadata !== tokenAccount) {
    return new Response('unauthorized', {status: 401})
  }

  //not found if value is null
  if (getWithMetadataResult.value === null) {
    return new Response('not found', {status: 404})
  }
  return new Response(getWithMetadataResult.value)
}

const worker: { fetch: (request: Request, env: Env) => Promise<Response> } = { fetch: handleRequest };
export default worker;


export function getToken(request: Request): string | null {
  const authorization = request.headers.get('Authorization')
  if (!authorization) {
    return null
  }
  const parts = authorization.split(' ')
  if (parts[0] !== 'Bearer') {
    return null
  }
  return parts[1]
}

export async function getTokenAccount(token: string, env: Env): Promise<string | null> {
  const accountId = await env.TOKEN.get(token)
  return accountId
}