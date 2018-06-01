package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type SomeChan struct {
	PathChan   chan string
	fallowChan chan string
}

//爬虫进程
type Process struct {
	rc chan string
	wc chan User
	pc chan string
	rp ReadProcess  // 读进程
	wp WriteProcess // 写进程
}

// 读进程
type ReadProcess interface {
	ReadContent(rc chan string, somechan *SomeChan)
}

// 写进程
type WriteProcess interface {
	WriteContent(wc chan User)
}

//通过url读取方式
type ReadUrlContent struct {
	base_url string
}

//写入mysql方式
type WriteToMysql struct {
	mysqldb string //"root:root@tcp('127.0.0.1:3306')/test"
}

//写入记事本
type WriteToNotebook struct {
	path string
}

//通过url获取内容
func (r *ReadUrlContent) ReadContent(rc chan string, somechan *SomeChan) {
	for {

		if url, ok := <-somechan.PathChan; ok {
			fmt.Println("start  ", url, time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05 "))
			resp, err := http.Get(url)
			if err != nil {
				log.Println(fmt.Sprint("read err", err.Error()))
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(fmt.Sprint("read err", err.Error()))
				continue
			}
			rc <- string(body)

		}
	}

}

//写入mysql
func (w *WriteToMysql) WriteContent(wc chan User) {
	fmt.Println(w.mysqldb)
	fmt.Println(<-wc)
}

//写入记事本
func (w *WriteToNotebook) WriteContent(wc chan User) {
	for {
		if user, ok := <-wc; ok {
			handel, err := os.OpenFile(w.path, os.O_APPEND, 7777)
			if err != nil {
				panic(fmt.Sprint(w.path+" file is not exist  info:", err.Error()))
			}
			defer handel.Close()

			str, _ := json.Marshal(user) //user结构中过的字段必须为可导出字段即字段首字母大写
			handel.WriteString(string(str) + "\n")
			fmt.Println("end  ", user.Name, time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05 "))

		}
	}

}

//解析数据放入结构体传到wc通道
func (b *Process) DoProcess(s *SomeChan) {
	for {
		if data, ok := <-b.rc; ok && true {
			var user User
			//messagehttps://my.csdn.net/qq_28602957
			// 查找连续的小写字母
			reg := regexp.MustCompile(`username='[\w]+'`)
			names := reg.FindAllString(data, -1)
			user.Follow = Fallow(names)

			reg = regexp.MustCompile(`<dd[\s]+class="person-detail"[\s]*>(.)*(\n)*(.)*</dd>`)
			persondetail := reg.FindAllString(data, -1)

			if persondetail == nil {
				user.Person_detail = ""
			} else {

				user.Person_detail = PersonDetail(string(persondetail[0]))
			}

			reg = regexp.MustCompile(`<dd[\s]+class="person-sign"[\s]*>(.)*(\n)*(.)*</dd>`)
			personsign := reg.FindAllString(data, -1)
			if personsign == nil {
				user.Person_sign = ""
			} else {

				user.Person_sign = PersonSign(string(personsign[0]))
			}

			reg = regexp.MustCompile(`<dt[\s]+class="person-nick-name"[\s]*>(.)*(\n)*(.)*</dt>`)
			namestr := reg.FindAllString(data, -1)
			if namestr != nil {
				user.Name = NameStr(string(namestr[0]))
			} else {
				user.Name = ""
			}

			b.wc <- user

			if user.Follow == "" {
				s.PathChan <- "https://my.csdn.net/qq_28602957"
				fmt.Println("process     >>", user.Name)
				continue
			}
			s.fallowChan <- user.Follow
			fmt.Println("process", user.Name, time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05 "))
		}
	}

}

func (s SomeChan) SetPathChan() {
	for {
		if follow, ok := <-s.fallowChan; ok && true {
			fmt.Println(follow)
			names := strings.Split(follow, ",")
			if len(names) == 1 {
				s.PathChan <- "https://my.csdn.net/" + names[0]
				fmt.Println("follow1  ", names[0], len(names))
				continue
			}
			if len(names) > 0 && follow != "" {
				num := 8 //关注长度有8个所以可以获取这个8个中的内容
				if len(names) < 8 {
					num = len(names)
				}
				a := rand.Intn(num)
				b := rand.Intn(num)
				for i, v := range names { //因为names中数量太多导致试用
					if i == a || i == b {
						s.PathChan <- "https://my.csdn.net/" + v
						fmt.Println("follow2  ", v, len(names), i)
					}
				}
			} else {
				s.PathChan <- "https://my.csdn.net/" + follow
			}
		}
	}

}

func Fallow(names []string) string {
	str := ""
	a := 0
	for _, v := range names {
		if a == 0 {
			str += string(v)
		} else {
			str += "," + string(v)
		}
		a++
	}
	str = strings.Replace(str, "username=", "", -1)
	str = strings.Replace(str, "'", "", -1)
	return str
}
func NameStr(str string) string {
	str = strings.Replace(str, "<", "", -1)
	str = strings.Replace(str, ">", "", -1)
	str = strings.Replace(str, "dt", "", -1)
	str = strings.Replace(str, "/", "", -1)
	str = strings.Replace(str, "class=\"person-nick-name\"", "", -1)
	str = strings.Replace(str, "span", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

func PersonSign(str string) string {
	str = strings.Replace(str, "<", "", -1)
	str = strings.Replace(str, ">", "", -1)
	str = strings.Replace(str, "/", "", -1)
	str = strings.Replace(str, "dd", "", -1)
	str = strings.Replace(str, "class=\"person-sign\"", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

func PersonDetail(str string) string {
	str = strings.Replace(str, "<", "", -1)
	str = strings.Replace(str, ">", "", -1)
	str = strings.Replace(str, "/", "", -1)
	str = strings.Replace(str, "dd", "", -1)
	str = strings.Replace(str, "class=\"person-detail\"", "", -1)
	str = strings.Replace(str, "span|span", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

//用户详情
type User struct {
	Username      string //用户名
	Name          string //用户名称
	Person_detail string //用户信息
	Person_sign   string //用户专业介绍
	Follow        string // 为username   关注列粉丝表以逗号隔开用户名
}

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano()) //加随机种子
	read := &ReadUrlContent{
		base_url: "https://my.csdn.net/qq_28602957",
	}
	write := &WriteToNotebook{
		path: "csdn.log",
	}
	somechan := &SomeChan{
		PathChan:   make(chan string, 100),
		fallowChan: make(chan string, 100),
	}
	p := &Process{
		rc: make(chan string, 100),
		wc: make(chan User, 100),
		rp: read,
		wp: write,
	}

	go func(somechan *SomeChan, url string) {
		somechan.PathChan <- url
	}(somechan, read.base_url)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go p.rp.ReadContent(p.rc, somechan)
	}
	wg.Add(2)
	go p.DoProcess(somechan)
	go somechan.SetPathChan()

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go p.wp.WriteContent(p.wc)
	}
	wg.Wait()

}

//`
//			<!DOCTYPE html>
//<html>
//  <head>
// </head>
//    <div class="main clearfix">
//      <div class="persional_property">
//        <div class="person_info_con">
//            <a class="person_add_focus"><i class="icon-plus"></i>关注</a>
//          <dl class="person-photo">
//            <dt><a href="javascript:;"><img src="https://avatar.csdn.net/8/D/5/1_stpeace.jpg?1527671342" class="header"></a></dt>
//            <dd class="focus_num"><b>31</b>关注</dd>
//            <dd class="fans_num"><b>2512</b>粉丝</dd>
//          </dl>

//        </div>
//      </div>
//      <div class="persion_section">
//        <div class="mod_resource">
//          <div id="tabResources" class="tabs clearfix">                  <a href="#" data-modal="tab" class="active">发表的博客</a>
//                  <a href="#" data-modal="tab">发布的帖子</a>
//              <a href="#" onclick="window.open('https://download.csdn.net/user/stpeace/uploads')" >贡献的资源</a></div>
//          <div id="divResources">

//            <div data-modal="tab-layer" class="list-blog list activeContent">
//                <input type="hidden" value="1" id="user_blog_page">
//              <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80383002" class="tit" title="520最悲情的告白是这样的">520最悲情的告白是这样的</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（2）</span>
//                        <span class="dTime">140人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80327276" class="tit" title="memory hierarchy------晒图">memory hierarchy------晒图</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（1）</span>
//                        <span class="dTime">120人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80368622" class="tit" title="相对路径究竟是相对谁的路径？">相对路径究竟是相对谁的路径？</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（1）</span>
//                        <span class="dTime">112人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80297997" class="tit" title="真正的出路：重读任正非2012实验室讲话------任总毕竟是任总啊">真正的出路：重读任正非2012实验室讲话------任总毕竟是任总啊</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（2）</span>
//                        <span class="dTime">459人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80295964" class="tit" title="腾讯专家分享：腾讯做业务监控的心得和经验">腾讯专家分享：腾讯做业务监控的心得和经验</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（0）</span>
//                        <span class="dTime">278人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80294283" class="tit" title="内无干货，慎入------负载均衡、过载保护、动态容错、名字服务">内无干货，慎入------负载均衡、过载保护、动态容错、名字服务</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（0）</span>
//                        <span class="dTime">221人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80279724" class="tit" title="makefile ifeq之坑: 1. syntax error near unexpected token  2.  *** missing separator.  Stop.">makefile ifeq之坑: 1. syntax error near unexpected token  2.  *** missing separator.  Stop.</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（0）</span>
//                        <span class="dTime">152人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80217315" class="tit" title="记一次Content-Length引发的血案">记一次Content-Length引发的血案</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（0）</span>
//                        <span class="dTime">243人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80037826" class="tit" title="图灵停机问题（The Halting Problem）------巧妙的证明">图灵停机问题（The Halting Problem）------巧妙的证明</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（0）</span>
//                        <span class="dTime">199人阅读</span>
//                    </li>
//                                    <li class="clearfix">
//                        <a href="http://blog.csdn.net/stpeace/article/details/80216124" class="tit" title="csrf攻击模拟">csrf攻击模拟</a>
//                        <span class="commend"><i class="comm_icon"></i>评论（0）</span>
//                        <span class="dTime">207人阅读</span>
//                    </li>

//                <div class="csdn-pagination">
//                    共2328个  &nbsp; 共233页&nbsp;<strong><a href="#" class="btn btn-xs btn-default active">1</a>&nbsp</strong>&nbsp;<a class="btn btn-xs btn-default" href="javascript:void(0);/2">2</a>&nbsp&nbsp;<a class="btn btn-xs btn-default" href="javascript:void(0);/3">3</a>&nbsp<span class="ellipsis">...</span>&nbsp;&nbsp;<a class="btn btn-xs btn-default" href="javascript:void(0);/233">233</a>&nbsp;<a class="btn btn-xs btn-default btn-next" href="javascript:void(0);/2">&gt;</a>&nbsp&nbsp;                </div>

//            </div>

//            <div data-modal="tab-layer" class="list-posts list">

//            </div>
//            <div data-modal="tab-layer" class="list-resource list ">

//            </div>
//          </div>
//        </div>
//        <div class="person_detail_tab2">
//          <ul id="tabDetail">
//            <li data-modal="tab" style="width: 100%" class="current_detail">最新动态</li>
//          </ul>
//        </div>
//        <div id="divDetail" class="aboutMe">
//          <div nodeType="myDetails" nodeIndex="my0" data-modal="tab-layer" class="myDetails">
//            <div class="mod_field_skill">
//              <div class="field">
//                <h3>熟悉的领域</h3>
//                <div class="tags clearfix">

//                </div>
//              </div>
//            </div>
//            <div class="mod_field_skill">
//              <div class="skill">
//                <h3>专业技能</h3>
//                <div class="tags clearfix">

//                </div>
//              </div>
//            </div>
//            <div class="person_education">
//              <h3><span>教育经历</span></h3>

//            </div>
//            <div class="person_job">
//              <h3><span>工作经历</span></h3>

//            </div>
//            <div class="mod_contact">
//              <h3>联系方式</h3>
//              <ul class="clearfix">
//                <li>电子邮箱：<span nodeType="email" class="email"></span></li>
//                <li>手机号码：<span nodeType="modile" class="modile"></span></li>
//                <li>QQ号码：<span nodeType="qq" class="qq"></span></li>
//                <li>微信号：<span nodeType="weixin" class="weixin"></span></li>
//              </ul>
//            </div>
//          </div>
//          <div nodeType="myNews" data-modal="tab-layer" class="myNews activeContent">
//            <div class="mod_per_dy">
//              <h3><a href="#" data-type="mine" class="mineTab active">我的全部动态<i class="icon"></i></a></h3>
//              <div class="mine">
//                <ul data-bind="my">
//                  <li><span data-bind="myTitle" class="info"></span>
//                    <div data-bind="myCont" class="cont"></div><span data-bind="myTime" class="time"></span>
//                  </li>
//                </ul>
//                <div class="more"><span class="icon-angle-down"></span><a href="#">显示更多</a></div>
//              </div>
//            </div>
//          </div>
//        </div>
//      </div>
//      <div class="persion_article">
//        <div class="interested_con">
//          <h3>对Ta感兴趣的人</h3>
//          <ul nodetype="inter-list" data-bind="list" class="clearfix">
//            <li><a href="#" target="_blank" data-bind="headerHref"><img src="" username="" data-bind="headerSrc"></a></li>
//          </ul>
//          <div class="count-and-more"><span>最近一周被访问了<em data-bind="times"></em>次</span></div>
//        </div>
//        <div class="mod_relations">
//          <h3>Ta的关系</h3>
//          <div class="list">
//            <div class="focus">
//              <div class="num">关注：<span>31</span>人</div>
//              <div class="header clearfix">
//              				<a href=broadview2006 ><img src='https://avatar.csdn.net/F/E/7/1_broadview2006.jpg?1527671342' username='broadview2006'/></a><a href=jj12345jj198999 ><img src='https://avatar.csdn.net/A/0/5/1_jj12345jj198999.jpg?1527671342' username='jj12345jj198999'/></a><a href=zhugeaming2018 ><img src='https://avatar.csdn.net/8/5/A/1_zhugeaming2018.jpg?1527671342' username='zhugeaming2018'/></a><a href=k346k346 ><img src='https://avatar.csdn.net/6/3/F/1_k346k346.jpg?1527671342' username='k346k346'/></a><a href=chenglinhust ><img src='https://avatar.csdn.net/4/2/A/1_chenglinhust.jpg?1527671342' username='chenglinhust'/></a><a href=zhouzxi ><img src='https://avatar.csdn.net/3/6/2/1_zhouzxi.jpg?1527671342' username='zhouzxi'/></a><a href=absurd ><img src='https://avatar.csdn.net/6/5/9/1_absurd.jpg?1527671342' username='absurd'/></a><a href=fullsail ><img src='https://avatar.csdn.net/1/5/E/1_fullsail.jpg?1527671342' username='fullsail'/></a>              </div>
//            </div>
//            <div class="focus beFocus">
//              <div class="num">被关注：<span>2512</span>人</div>
//              <div class="header clearfix">
//              					<a href=qq_41499763 ><img src='https://avatar.csdn.net/C/5/B/1_qq_41499763.jpg?1527671342' username='qq_41499763'/></a><a href=sty945 ><img src='https://avatar.csdn.net/6/2/8/1_sty945.jpg?1527671342' username='sty945'/></a><a href=qq_37583486 ><img src='https://avatar.csdn.net/1/5/A/1_qq_37583486.jpg?1527671342' username='qq_37583486'/></a><a href=alisa_xf ><img src='https://avatar.csdn.net/D/7/5/1_alisa_xf.jpg?1527671342' username='alisa_xf'/></a><a href=u013258199 ><img src='https://avatar.csdn.net/B/A/4/1_u013258199.jpg?1527671342' username='u013258199'/></a><a href=dahe8610 ><img src='https://avatar.csdn.net/F/A/B/1_dahe8610.jpg?1527671342' username='dahe8610'/></a><a href=cndm123 ><img src='https://avatar.csdn.net/9/3/1/1_cndm123.jpg?1527671342' username='cndm123'/></a><a href=zx_water ><img src='https://avatar.csdn.net/5/3/E/1_zx_water.jpg?1527671342' username='zx_water'/></a>
//              </div>
//            </div>
//          </div>
//        </div>
//            </div>
//    </div>

//  </body>
//</html>
//			`
