package main

import (
	"fmt"
	"html/template"
	"os"
)

type fluentd struct {
	Path string
	Tag  string
	Name string
}

func main() {
	tmpStr := `
<source>
  @type tail
  path {{.Path}} # 要收集的日志文件路径
  pos_file {{.Path}}.pos
  tag {{.Tag}} # 标签
  # 分析提取日志信息,最终的正则表达式以日志的字段为准。
  <parse>
    @type regexp
    expression /^\[(?<logtime>[^\]]*)\]\s+\[(?<logtype>[^ ]*)\]\s+(?<loglevel>[^ ]*)\s+(?<logcontent>sysLog\s+.*)/
    # time_key logtime
    time_format %Y-%m-%d %H:%M:%s
  </parse>
</source>

<match {{.Tag}}>
  @type copy
  <store>
    # debug模式 (tail -f /var/log/td-agent/td-agent.log)可以查看日志收集情况
    @type stdout
  </store>
  <store>
    @type elasticsearch # 使用的数据存储的插件
    host es.example.com # es服务器地址
    port 31114      # es端口
    logstash_format false
    flush_interval 10
    index_name {{.Name}}
  </store>
</match>
`
	hisp := fluentd{
		Path: "/glzt/his_new/hisp_resource/logs/hispRes/sys.log",
		Tag:  "hisp_resource.log",
		Name: "br-his-hisp-resource",
	}
	schedule := fluentd{
		Path: "/glzt/his_new/schedule/logs/schedule/sys.log",
		Tag:  "schedule.log",
		Name: "br-his-schedule",
	}

	socketserver := fluentd{
		Path: "/glzt/his_new/socketserver/logs/socketserver/sys.log",
		Tag:  "socketserver.log",
		Name: "br-his-socketserver",
	}

	outp := fluentd{
		Path: "/glzt/his_new/outp_doctor/logs/outpdr/sys.log",
		Tag:  "outp_doctor.log",
		Name: "br-his-outp-doctor",
	}

	templateEngine := fluentd{
		Path: "/glzt/his_new/template_engine/logs/templateengine/sys.log",
		Tag:  "template_engine.log",
		Name: "br-his-template-engine",
	}

	task := fluentd{
		Path: "/glzt/his_new/task/logs/task/sys.log",
		Tag:  "task.log",
		Name: "br-his-task",
	}

	medinsurance := fluentd{
		Path: "/glzt/his_new/medinsurance/logs/medinsurance/sys.log",
		Tag:  "medinsurance.log",
		Name: "br-his-medinsurance",
	}

	dictManage := fluentd{
		Path: "/glzt/his_new/dict_manage/logs/dictManage/sys.log",
		Tag:  "dict_manage.log",
		Name: "br-his-dict-manage",
	}

	userManage := fluentd{
		Path: "/glzt/his_new/user_manager/logs/usermanage/sys.log",
		Tag:  "user_manager.log",
		Name: "br-his-user-manager",
	}

	fileManage := fluentd{
		Path: "/glzt/his_new/file_manage/logs/filemanage/sys.log",
		Tag:  "file_manage.log",
		Name: "br-his-file-manage",
	}

	printManage := fluentd{
		Path: "/glzt/his_new/print_manage/logs/printmanage/sys.log",
		Tag:  "print_manage.log",
		Name: "br-his-print-manage",
	}

	register := fluentd{
		Path: "/glzt/his_new/register/logs/register/sys.log",
		Tag:  "register.log",
		Name: "br-his-register",
	}

	hisadapter := fluentd{
		Path: "/glzt/his_new/hisadapter/logs/hisadapter/sys.log",
		Tag:  "hisadapter.log",
		Name: "br-his-hisadapter",
	}

	orgmanage := fluentd{
		Path: "/glzt/his_new/org_manager/logs/orgmanage/sys.log",
		Tag:  "orgmanage.log",
		Name: "br-his-orgmanage",
	}

	appoint := fluentd{
		Path: "/glzt/his_new/appoint/logs/appoint/sys.log",
		Tag:  "appoint.log",
		Name: "br-his-appoint",
	}

	fluentds := make([]fluentd, 0)
	fluentds = append(fluentds, hisp, schedule, socketserver, outp, templateEngine, task, medinsurance, dictManage, userManage, fileManage, printManage, register, hisadapter, orgmanage, appoint)

	//fluentdsJson, err := json.Marshal(&fluentds)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//fmt.Printf("%v\n", fluentdsJson)

	for _, f := range fluentds {
		fmt.Printf("%v\n", f)
		temp := template.New("test")
		temp, err := temp.Parse(tmpStr)
		if err != nil {
			panic(err)
		}
		err = temp.Execute(os.Stdout, &f)
		if err != nil {
			panic(err)
		}
	}

}
