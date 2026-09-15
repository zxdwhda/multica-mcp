#!/usr/bin/env python3
"""Real OAuth + MCP smoke. Tokens stay in memory. Optional isolated write/cleanup test."""
import argparse,base64,hashlib,json,pathlib,re,secrets,urllib.request,urllib.parse,urllib.error
class NoRedirect(urllib.request.HTTPRedirectHandler):
 def redirect_request(self,*args,**kwargs):return None
opener=urllib.request.build_opener(NoRedirect)
def main():
 parser=argparse.ArgumentParser();parser.add_argument("--base",default="https://mcp.wildflow.cn");parser.add_argument("--write",action="store_true");parser.add_argument("--full-api",action="store_true");parser.add_argument("--redirect-uri",default="https://chatgpt.com/connector_platform/oauth/callback");args=parser.parse_args()
 base=args.base.rstrip('/');origin="https://mcp.wildflow.cn";resource=origin+"/multica/mcp"
 cfg=json.loads((pathlib.Path.home()/".multica/config.json").read_text());checks=[]
 def call(path,body=None,headers=None,method=None):
  h=headers or {};r=urllib.request.Request(base+path,data=body,headers=h,method=method)
  try:resp=opener.open(r,timeout=45)
  except urllib.error.HTTPError as e:resp=e
  raw=resp.read();return resp.status,resp.headers,raw
 def expect(label,status,want):
  assert status==want,f"{label}: HTTP {status}, expected {want}"
  checks.append({"check":label,"status":status});print(label,status,flush=True)
 def form(path,data,headers=None):return call(path,urllib.parse.urlencode(data).encode(),{"Content-Type":"application/x-www-form-urlencoded",**(headers or {})})
 status,h,raw=call('/multica/healthz');expect('health',status,200)
 status,h,raw=call('/multica/mcp',b'{}',{'Content-Type':'application/json'});expect('unauthenticated MCP',status,401);assert 'resource_metadata' in h.get('WWW-Authenticate','')
 status,h,raw=call('/multica/register',json.dumps({'redirect_uris':[args.redirect_uri],'token_endpoint_auth_method':'none'}).encode(),{'Content-Type':'application/json'});expect('DCR',status,201);client=json.loads(raw)['client_id']
 verifier=secrets.token_urlsafe(48);challenge=base64.urlsafe_b64encode(hashlib.sha256(verifier.encode()).digest()).decode().rstrip('=')
 q={'client_id':client,'redirect_uri':args.redirect_uri,'response_type':'code','code_challenge_method':'S256','code_challenge':challenge,'resource':resource,'state':'smoke-'+secrets.token_hex(8),'scope':'multica:access'}
 status,h,raw=call('/multica/authorize?'+urllib.parse.urlencode(q));expect('consent page',status,200);cookie=h['Set-Cookie'].split(';',1)[0];pending=re.search(r'name="request" value="([^"]+)"',raw.decode()).group(1)
 status,h,raw=form('/multica/authorize',{'request':pending,'pat':cfg['token']},{'Origin':origin,'Cookie':cookie});
 if status!=303:print('approval error',raw.decode()[:400],dict((k,v) for k,v in h.items() if k.lower() in ['x-fc-error-type','content-type']),flush=True)
 expect('PAT approval',status,303);redirect=urllib.parse.urlparse(h['Location']);params=urllib.parse.parse_qs(redirect.query);assert params['state'][0]==q['state'];code=params['code'][0]
 f={'grant_type':'authorization_code','client_id':client,'code':code,'code_verifier':verifier,'redirect_uri':q['redirect_uri'],'resource':resource}
 status,h,raw=form('/multica/token',f);expect('code exchange',status,200);tokens=json.loads(raw);status,h,raw=form('/multica/token',f);expect('code replay rejected',status,400)
 access=tokens['access_token'];count=0
 def rpc(method,params):
  nonlocal count;count+=1
  status,h,raw=call('/multica/mcp',json.dumps({'jsonrpc':'2.0','id':count,'method':method,'params':params}).encode(),{'Authorization':'Bearer '+access,'Content-Type':'application/json','Accept':'application/json, text/event-stream','MCP-Protocol-Version':'2025-06-18'})
  expect(method if method!='tools/call' else params['name'],status,200);data=json.loads(raw);assert 'error' not in data,data.get('error');return data['result']
 def tool(name,arguments):
  data=rpc('tools/call',{'name':name,'arguments':arguments});assert not data.get('isError'),data.get('content');return data.get('structuredContent',{})
 rpc('initialize',{'protocolVersion':'2025-06-18','capabilities':{},'clientInfo':{'name':'deployment-smoke','version':'1'}})
 catalog=rpc('tools/list',{});api_catalog=json.loads((pathlib.Path(__file__).resolve().parents[1]/'internal/apicatalog/catalog.json').read_text());assert len(catalog['tools'])==17+len(api_catalog['operations'])
 for name in ['multica_list_projects','multica_list_tasks','multica_list_agents','multica_list_statuses']:tool(name,{})
 if args.full_api:
  for name in ['list_projects','list_issues','list_agents','list_labels','list_skills','list_squads','list_autopilots','list_agent_runtimes','list_workspaces']:
   result=tool('multica_api_'+name,{})
   assert result['status']==200,(name,result['status'])
  marker='mcp-api-verification-'+secrets.token_hex(4)
  for kind,body,field,updated in [
   ('project',{'title':marker,'description':'Temporary API verification','status':'planned'},'description','Updated through full API'),
   ('label',{'name':marker,'color':'#3B82F6','description':'Temporary API verification','resource_type':'issue'},'description','Updated through full API'),
   ('skill',{'name':marker,'description':'Temporary API verification','content':'# Verification\nTemporary test only.'},'description','Updated through full API')]:
   created_id=None
   try:
    created_result=tool('multica_api_create_'+kind,{'body':body})
    assert created_result['status'] in [200,201],created_result['status']
    created_body=created_result['body'];created_id=created_body['id']
    changed=tool('multica_api_update_'+kind,{'path':{'id':created_id},'body':{field:updated}})
    assert changed['status']==200,changed['status']
    if kind!='label':
     got=tool('multica_api_get_'+kind,{'path':{'id':created_id}})
     assert got['body'][field]==updated
    else:
     got=tool('multica_api_list_labels',{})
     rows=got['body'] if isinstance(got['body'],list) else got['body']['labels']
     assert any(x['id']==created_id and x[field]==updated for x in rows)
   finally:
    if created_id:
     removed=tool('multica_api_delete_'+kind,{'path':{'id':created_id}})
     assert removed['status'] in [200,204],removed['status']
 if args.write:
  created=[]
  try:
   marker='MCP deployment smoke '+secrets.token_hex(4)
   task=tool('multica_create_task',{'title':marker,'description':'Temporary integration test, unassigned.','status':'backlog'});created.append(task['id'])
   tool('multica_update_task',{'task_id':task['id'],'description':'','suppress_run':True})
   detail=tool('multica_get_task',{'task_id':task['id']});assert detail.get('description') in ['',None];assert not detail.get('warnings'),detail.get('warnings')
   tool('multica_add_comment',{'task_id':task['id'],'comment':'/note MCP deployment verification; no agent should run.'})
   tool('multica_search_tasks',{'query':marker,'limit':10})
  finally:
   for id in reversed(created):
    r=urllib.request.Request(cfg['server_url']+'/api/issues/'+id,method='DELETE',headers={'Authorization':'Bearer '+cfg['token'],'X-Workspace-ID':cfg['workspace_id']})
    with urllib.request.urlopen(r,timeout=30) as resp:expect('cleanup temporary issue',resp.status,204)
 status,h,raw=form('/multica/token',{'grant_type':'refresh_token','client_id':client,'refresh_token':tokens['refresh_token'],'resource':resource});expect('refresh token',status,200);new=json.loads(raw);access=new['access_token'];tool('multica_list_projects',{})
 status,h,raw=form('/multica/revoke',{'client_id':client,'token':new['refresh_token']});expect('revoke',status,200)
 status,h,raw=call('/multica/mcp',b'{}',{'Authorization':'Bearer '+access,'Content-Type':'application/json'});expect('revoked access rejected',status,401)
 out=pathlib.Path(__file__).resolve().parents[1]/'deploy/build';out.mkdir(parents=True,exist_ok=True);(out/'smoke-results.json').write_text(json.dumps({'base':base,'checks':checks,'tools':len(catalog['tools']),'write_test':args.write,'full_api_test':args.full_api},indent=2));print('All smoke checks passed',flush=True)
if __name__=='__main__':main()
