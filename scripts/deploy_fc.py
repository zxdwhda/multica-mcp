#!/usr/bin/env python3
"""Deploy only the Multica function via Aliyun CLI. Never modifies shared DNS/domain/OSS."""
import base64, json, os, pathlib, subprocess, tempfile, zipfile
ROOT=pathlib.Path(__file__).resolve().parents[1]
REGION="ap-southeast-1"
NAME="multica-mcp"
def cli(args,body=None,allow_missing=False):
    path=None
    try:
        if body is not None:
            fd,path=tempfile.mkstemp(prefix="multica-fc-",suffix=".json")
            with os.fdopen(fd,"w") as f: json.dump(body,f)
            args=args+["--body-file",path]
        p=subprocess.run(["aliyun"]+args,capture_output=True,text=True)
        if p.returncode:
            if allow_missing and ("FunctionNotFound" in p.stderr or "NotFound" in p.stderr):return None
            # Avoid printing echoed request bodies or environment values.
            raise RuntimeError("Aliyun CLI failed: "+p.stderr[:400])
        return json.loads(p.stdout)
    finally:
        if path: os.unlink(path)
def main():
    if subprocess.check_output(["git","status","--porcelain"],cwd=ROOT,text=True).strip():
        raise RuntimeError("Commit changes before deployment")
    revision=subprocess.check_output(["git","rev-parse","HEAD"],cwd=ROOT,text=True).strip()
    cfg=json.loads((pathlib.Path.home()/".multica/config.json").read_text())
    if cfg["server_url"]!="https://api.multica.ai":raise RuntimeError("Unexpected Multica API target")
    role=cli(["ram","GetRole","--RoleName","AliyunFcDefaultRole"])["Role"]["Arn"]
    build=ROOT/"deploy/build";build.mkdir(parents=True,exist_ok=True)
    env={**os.environ,"GOOS":"linux","GOARCH":"amd64","CGO_ENABLED":"0"}
    subprocess.run(["go","build","-trimpath","-ldflags=-s -w -X multica-mcp/internal/version.Revision="+revision,"-o",str(build/"bootstrap"),"."],cwd=ROOT,env=env,check=True)
    with zipfile.ZipFile(build/"function.zip","w",zipfile.ZIP_DEFLATED) as z:z.write(build/"bootstrap","bootstrap")
    variables={"MCP_TRANSPORT":"http","MCP_HTTP_PORT":"9000","MCP_HTTP_PREFIX":"/multica","MCP_OAUTH_ORIGIN":"https://mcp.wildflow.cn","MCP_OSS_BUCKET":"wildflow-mcp-state-sg","MCP_OSS_ENDPOINT":"https://oss-ap-southeast-1-internal.aliyuncs.com","MULTICA_BASE_URL":cfg["server_url"],"MULTICA_TOKEN":cfg["token"],"MULTICA_WORKSPACE_ID":cfg["workspace_id"],"MULTICA_READ_ONLY":"false","LOG_LEVEL":"info"}
    body={"description":"Multica MCP "+revision,"logConfig":{"project":"wildflow-mcp-sg","logstore":"multica-mcp","enableRequestMetrics":True,"enableInstanceMetrics":True,"logBeginRule":"None"},"runtime":"custom.debian11","handler":"bootstrap","cpu":0.2,"memorySize":256,"diskSize":512,"timeout":120,"instanceConcurrency":20,"internetAccess":True,"role":role,"environmentVariables":variables,"customRuntimeConfig":{"command":["/code/bootstrap"],"port":9000,"healthCheckConfig":{"httpGetUrl":"/multica/healthz","initialDelaySeconds":0,"periodSeconds":3,"timeoutSeconds":2,"failureThreshold":3,"successThreshold":1}},"code":{"zipFile":base64.b64encode((build/"function.zip").read_bytes()).decode()}}
    base="/2023-03-30/functions"
    old=cli(["fc","GET",base+"/"+NAME,"--region",REGION],allow_missing=True)
    if old is None:
        body["functionName"]=NAME
        result=cli(["fc","POST",base,"--region",REGION],body)
    else:
        # Preserve unrelated environment keys and never target other functions.
        body["environmentVariables"]={**old.get("environmentVariables",{}),**variables}
        result=cli(["fc","PUT",base+"/"+NAME,"--region",REGION],body)
    triggers=cli(["fc","GET",base+"/"+NAME+"/triggers","--region",REGION]).get("triggers",[])
    if not any(t["triggerName"]=="http" for t in triggers):
        cli(["fc","POST",base+"/"+NAME+"/triggers","--region",REGION],{"triggerName":"http","triggerType":"http","qualifier":"LATEST","triggerConfig":json.dumps({"authType":"anonymous","methods":["GET","POST","DELETE","OPTIONS"],"disableURLInternet":False})})
    result=cli(["fc","GET",base+"/"+NAME,"--region",REGION])
    triggers=cli(["fc","GET",base+"/"+NAME+"/triggers","--region",REGION])
    report={k:result.get(k) for k in ["functionName","runtime","codeChecksum","lastModifiedTime","state"]};report.update({"region":REGION,"qualifier":"LATEST","publicEndpoint":"https://mcp.wildflow.cn/multica/mcp","triggers":triggers.get("triggers",[])})
    (build/"deployment.json").write_text(json.dumps(report,indent=2))
    print(json.dumps(report,indent=2))
if __name__=="__main__":main()
