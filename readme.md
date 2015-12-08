#GOLANG搭建单、双向自认证HTTPS服务器
###前言
>2015年双11期间淘宝、天猫实现了全站式https安全传输，web安全问题已经成了人们关注的话题，那什么是https呢？如何实现单、双向自认证https服务器呢？接下来我们将一一介绍。

##一、HTTPS相关概念已经认证流程
######基本概念：
[**HTTPS**](http://baike.baidu.com/link?url=XuEFqp8HTAIWBO12QMzj54K1iIBGPL6VJGPEn85nyCirdG8LE104hMYvOeDgfucyMf3gu1zPLap3i0BKb-SKHa)（全称：Hyper Text Transfer Protocol over Secure Socket Layer），是以安全为目标的HTTP通道，简单讲是HTTP的安全版。即HTTP下加入SSL层，HTTPS的安全基础是SSL，因此加密的详细内容就需要SSL。 它是一个URI scheme（抽象标识符体系），句法类同http:体系。用于安全的HTTP数据传输。https:URL表明它使用了HTTP，但HTTPS存在不同于HTTP的默认端口及一个加密/身份验证层（在HTTP与TCP之间）。这个系统的最初研发由网景公司(Netscape)进行，并内置于其浏览器Netscape Navigator中，提供了身份验证与加密通讯方法。现在它被广泛用于万维网上安全敏感的通讯，例如交易支付方面。关于https详细介绍请见：[大型网站的HTTPS实践](http://studygolang.com/articles/2984)。

**SSL**(Secure Socket Layer)：是Netscape公司设计的主要用于WEB的安全传输协议。从名字就可以看出它在https协议栈中负责实现上面提到的加密层。

**数字证书**：一种文件的名称，好比一个机构或人的签名，能够证明这个机构或人的真实性。其中包含的信息，用于实现上述功能。

**加密和认证**：加密是指通信双方为了防止铭感信息在信道上被第三方窃听而泄漏，将明文通过加密变成密文，如果第三方无法解密的话，就算他获得密文也无能为力；认证是指通信双方为了确认对方是值得信任的消息发送或接受方，而不是使用假身份的非法者，采取的确认身份的方式。只有同时进行了加密和认证才能保证通信的安全，因此在SSL通信协议中这两者都被应。早期一般是用对称加密算法，现在一般都是不对称加密，最常见的算法就是RSA。

**消息摘要**：这个技术主要是为了避免消息被篡改。消息摘要是把一段信息，通过某种算法，得出一串字符串。这个字符串就是消息的摘要。如果消息被篡改（发生了变化），那么摘要也一定会发生变化（如果2个不同的消息生成的摘要是一样的，那么这就叫发生了碰撞）。消息摘要的算法主要有MD5和SHA，在证书领域，一般都是用SHA（安全哈希算法）。

数字证书、加密和认证、消息摘要三个技术结合起来，就是在HTTPS中广泛应用的证书（certificate），证书本身携带了加密/解密的信息，并且可以标识自己的身份，也自带消息摘要。
######HTTPS认证过程：
1. 浏览器发送一个连接请求给安全服务器。
2. 服务器将自己的证书，以及同证书相关的信息发送给客户浏览器。
3. 客户浏览器检查服务器送过来的证书是否是由自己信赖的 CA 中心所签发的。如果是，就继续执行协议；如果不是，客户浏览器就给客户一个警告消息：警告客户这个证书不是可以信赖的，询问客户是否需要继续。
4. 接着客户浏览器比较证书里的消息，例如域名和公钥，与服务器刚刚发送的相关消息是否一致，如果是一致的，客户浏览器认可这个服务器的合法身份。
5. 服务器要求客户发送客户自己的证书。收到后，服务器验证客户的证书，如果没有通过验证，拒绝连接；如果通过验证，服务器获得用户的公钥。
6. 客户浏览器告诉服务器自己所能够支持的通讯对称密码方案。
7. 服务器从客户发送过来的密码方案中，选择一种加密程度最高的密码方案，用客户的公钥加过密后通知浏览器。
8. 浏览器针对这个密码方案，选择一个通话密钥，接着用服务器的公钥加过密后发送给服务器。
9. 服务器接收到浏览器送过来的消息，用自己的私钥解密，获得通话密钥。
10. 服务器、浏览器接下来的通讯都是用对称密码方案，对称密钥是加过密的。

上面所述的是双向认证 SSL 协议的具体通讯过程，这种情况要求服务器和用户双方都有证书。单向认证 SSL 协议不需要客户拥有 CA 证书，具体的过程相对于上面的步骤，只需将服务器端验证客户证书的过程去掉，以及在协商对称密码方案，对称通话密钥时，服务器发送给客户的是没有加过密的 （这并不影响 SSL 过程的安全性）密码方案。这样，双方具体的通讯内容，就是加过密的数据，如果有第三方攻击，获得的只是加密的数据，第三方要获得有用的信息，就需要对加密 的数据进行解密，这时候的安全就依赖于密码方案的安全。而幸运的是，目前所用的密码方案，只要通讯密钥长度足够的长，就足够的安全。这也是我们强调要求使用128 位加密通讯的原因。
##二、自认证根证书
1. 创建根证书密钥文件(自己做CA)root.key：

	```shell
	$openssl genrsa -des3 -out root.key 2048
	```

	需要输入两次私钥密码            
	![](/res/1.png)
2. 创建根证书的申请文件root.csr：

	```shell
	$openssl req -new -key root.key -out root.csr
	```

	输入root.key的密码                             
	![](/res/2.png)

3. 创建根证书root.crt：
	
	```shell
	$openssl x509 -req -days 3650 -sha256 -extensions v3_ca -signkey root.key -in root.csr -out root.crt
	```

	生成根证书                                                        
	![](/res/3.png)
	
##三、SSL单向认证
1. 创建服务器证书秘钥

	```shell
	$openssl genrsa –des3 -out server.key 2048
	```

	需要输入两次私钥密码                                    
	![](/res/4.png)
2. 去除key口令
	
	```shell
	$openssl rsa -in server.key -out server.key
	```

	需要输入私钥密码                                           
	![](/res/5.png)
3. 创建服务器证书申请文件server.csr

	```shell
	$openssl req -new -key server.key -out server.csr
	```

	"Common Name"最好跟网站的域名一致                           
	![](/res/6.png)
4. 创建服务器证书server.crt

	```shell
	$openssl x509 -req -days 365 -sha256 -extensions v3_req -CA root.crt -CAkey root.key -CAcreateserial -in server.csr -out server.crt
	```

	需要输入根私钥密码                                     
	![](/res/7.png)                                          

5. 客户端导入根证书并添加到“信任的根服务站点”                                     
	![](/res/8.png)                                   
	![](/res/9.png)                               
	![](/res/10.png)                       
	![](/res/11.png)                            
	![](/res/12.png)                           
	![](/res/13.png)                                 
	![](/res/14.png)                                   
	![](/res/15.png)                                         
6. golang实现简单的https服务器
	
	```Go
	package main

	import (
		"io"
		"log"
		"net/http"
	)
	
	func handler(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "golang https server")
	}
	
	func main() {
		http.HandleFunc("/", handler)
		if err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}

	```

7. 在浏览器中测试                                             
	![](/res/16.png)	

##四、SSL双向认证
在单向认证的基础上添加客户端证书并在golang服务器源码上添加客户端认证相关代码

1. 创建客户端证书私钥

	```shell
	$openssl genrsa -des3 -out client.key 2048
	```	

	需要输入两次私钥密码                                       
	![](/res/17.png)
2. 去除key口令
	
	```shell
	$openssl rsa -in client.key -out client.key
	```

	需要输入私钥密码                                           
	![](/res/18.png)
3. 创建客户端证书申请文件client.csr

	```shell
	$openssl req -new -key client.key -out client.csr
	```

	![](/res/19.png)                                      

3. 创建客户端证书文件client.crt

	```shell
	$openssl x509 -req -days 365 -sha256 -extensions v3_req -CA root.crt -CAkey root.key -CAcreateserial -in client.csr -out client.crt
	```

	![](/res/20.png)                      
4. 将客户端证书文件client.crt和客户端证书密钥文件client.key合并成客户端证书安装包client.pfx

	```shell
	$openssl pkcs12 -export -in client.crt -inkey client.key -out client.pfx
	```	

	设置客户端安装时的密码                                           
	![](/res/23.png)
2. 添加客户端证书

	参见服务器端添加证书，客户端证书添加到“个人”里面就可以                      
	![](/res/21.png)
3. 修改服务器代码
	
	```Go
	package main
	
	import (
		"crypto/tls"
		"crypto/x509"
		"io"
		"io/ioutil"
		"log"
		"net/http"
	)
	
	type httpsHandler struct {
	}
	
	func (*httpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "golang https server!!!")
	}
	
	func main() {
		pool := x509.NewCertPool()
		caCertPath := "root.crt"
	
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			log.Fatal("ReadFile err:", err)
			return
		}
		pool.AppendCertsFromPEM(caCrt)
	
		s := &http.Server{
			Addr:    ":8080",
			Handler: &httpsHandler{},
			TLSConfig: &tls.Config{
				ClientCAs:  pool,
				ClientAuth: tls.RequireAndVerifyClientCert,
			},
		}
	
		if err = s.ListenAndServeTLS("server.crt", "server.key"); err != nil {
			log.Fatal("ListenAndServeTLS err:", err)
		}
	}

	```
4. 在浏览器中测试                                                
	![](/res/22.png)
5. 使用golang访问https服务器

	```Go
	package main

	import (
		"crypto/tls"
		"crypto/x509"
		"io/ioutil"
		"log"
		"net/http"
	)
	
	func main() {
		pool := x509.NewCertPool()
		caCertPath := "root.crt"
	
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			log.Fatal("ReadFile err:", err)
			return
		}
		pool.AppendCertsFromPEM(caCrt)
	
		cliCrt, err := tls.LoadX509KeyPair("client.crt", "client.key")
		if err != nil {
			log.Fatal("LoadX509KeyPair err:", err)
			return
		}
	
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get("https://localhost:8080")
		if err != nil {
			log.Fatal("client error:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
	}

	```	

###结语
希望通过这次实例能让大家更好的理解、应用https，谢谢观看。
