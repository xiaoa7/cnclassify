package main

import (
	cnc "cnclassify"
	"log"
	"net/http"
	"strings"
	"time"
)

//
var classify cnc.Classify

//
func init() {
	classify = cnc.Classify{Name: "测试分类"}
}

//性能测试的例子
func TestPerf() {
	//加载规则
	l1 := time.Now()
	classify.LoadRulesByDir("rules")
	log.Printf("加载规则共耗时:%.5f秒\n", time.Now().Sub(l1).Seconds())
	l1 = time.Now()
	for i := 0; i < 10000; i++ {
		text := "河南生政府采购网关于XX招标的中标公示"
		result := classify.Classification(text)
		log.Println(text, result)
	}
	log.Printf("分类10000次，共耗时:%.3f秒\n", time.Now().Sub(l1).Seconds())
}

//
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(`
		<!DOCTYPE HTML>
<HTML>
 <HEAD>
  <TITLE>中文信息分类识别</TITLE>
  <meta charset="utf-8">
  <META NAME="Generator" CONTENT="EditPlus">
  <META NAME="Author" CONTENT="">
  <META NAME="Keywords" CONTENT="">
  <META NAME="Description" CONTENT="">
  <script src="http://www.qimingxing.info/js/jquery.js"></script>
  <style>
   .text{height:24px;font-size:14px;padding:2px;}
   body{line-height:24px;}
   .submit{width:90px;height:32px;font-size:18px;font-weight:bold;}
  </style>
 </HEAD>
 <BODY>
 <h2>  中文分类系统 </h2>
 规则介绍：<br>
 示例:招标+(公示|公告|通知)<br/>
 解释:识别信息中有招标并且有（公示或公告或通知）的，认为满足此分类<br/>
 语法:支持以下符号+与的关系，|或的关系，^不包含的关系，可以用（）小括号，组织更复杂的表达式。<br/>
 如:招标+(公告|公示|通知)^(中标+(结果|公示))^(作废|流标|无效)<br/>
 <hr/>
 规则名称:<input type="text" id="name" class="text" value="招标" size="20"/><br/>
 规则:<input type="text" id="rule" class="text" value="招标+(公示|公告|通知)" size="120"/><br/>
 待测试的文本:<br/><textarea id="text" rows="10" cols="32"></textarea><br/>
 <input type="button" class="submit" onclick="doSubmit()" value="提交测试"/> 
 <SCRIPT LANGUAGE="JavaScript">
 <!--
	function doSubmit(){
		$.post("/",{name:$("#name").val(),rule:$("#rule").val(),text:$("#text").val()},function(data){
			alert("分类结果："+data);
		});
	}
 //-->
 </SCRIPT>
 </BODY>
</HTML>

		
		`))
	} else if r.Method == "POST" {
		r.ParseForm()
		name, rule, text := r.FormValue("name"), r.FormValue("rule"), r.FormValue("text")
		classify.LoadRulesByString(name, rule)
		tmp := classify.Classification(text)
		w.Write([]byte(strings.Join(tmp, ";")))
	}
}

//WEB端规则测试
func main() {
	http.HandleFunc("/", IndexHandle)
	http.ListenAndServe(":80", nil)
}
