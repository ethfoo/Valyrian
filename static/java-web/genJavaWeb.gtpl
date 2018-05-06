<html>
    <head></head>
    <body>
        <div><h2>生成java-web项目的shell脚本</h2></div>
        <div>
         <form action="/java-web/gen" method="POST">
            <h3>项目名称(必填)：</h3>
            <input type="text" name="shellName" />
            <h3>默认git地址(必填)</h3>
            <input type="text" name="gitAddr" />
            <input type="radio" name="useSsh" value="true" checked>git使用ssh方式</input>
            <input type="radio" name="useSsh" value="false">git使用https方式</input>
            <h3>git ssh权限：(使用ssh方式填写)</h3>
            私钥:<textarea style='width:200px,height:100px' name="sshPri"></textarea>
            known_hosts:<textarea style='width:200px,height:100px' name="sshHost"></textarea>
            <h3>git 用户名密码:（使用https的方式时填写）</h3>
            username:<input type="text" name="gitUsername" />
            password:<input type="password" name="gitPassword" />
            <h3>maven的settings文件，不传则使用默认maven settings</h3>
            <textarea style='width:200px,height:100px' name="mavenSettings"></textarea>
            <h3>是否有maven子模块</h3>
            <input type="radio" name="standalone" value="-s" checked>没有maven子模块</input>
            <input type="radio" name="standalone" value="">有maven子模块</input>
            
            <h3>生成的镜像Repository前缀，如hub.c.163.com/ncerepo/</h3>
            <input type="text" name="repository" value="{{.Repo}}"/>
            <input type="submit" value="生成脚本" />
        </form>
        </div>
        <div>
            
        </div>
        <div>
            
        </div>
        <div>
            
        </div>
        
    </body>
</html>