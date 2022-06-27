
export type Env = {
  FLAG: KVNamespace
}

export async function handleRequest(request: Request, env: Env) {
  const url = new URL(request.url);
  const projectId = url.pathname.slice(1)
  if (projectId.length !== 36) {
    return new Response('invalid project id', {status: 400})
  }
  const config = await env.FLAG.get(projectId)

  if (config === null) {
    return new Response('not found', {status: 404})
  }
  return new Response(config)
}

const worker: { fetch: (request: Request, env: Env) => Promise<void> } = { fetch: handleRequest };
export default worker;
