<html>
    <head></head>
    <body>
        <div><h2>运行shell脚本: {{.ShellName}}</h2></div>
        <div>
        <form action="/java-web/build" method="POST">

            <h3>默认git地址：</h3>
            <h3>git branch:</h3>
            <input type="text" name="branch" />
            <input type="text" name="ns" value={{.ShellName}} hidden />
            <h3>模块名称</h3>
            <input type="text" name="targetName" /> 
            <br />    
            <input type="submit" value="运行脚本" />
            <h3>-------------覆盖脚本默认参数----------</h3>
            
        </form>
        </div>
    </body>
</html>