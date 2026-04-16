package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// глобальные
var Keys = make(chan string)
var Inchat = false
var Conn2 net.Conn

//var CONNCALLs = false

//	type Сonnname struct {
//		name string
//		ip   string
//		port string
//	}
//
// поток для ожидания принятия запроса
func accepter(L net.Listener, txt *widget.Label) {
	//бесконечно слушает
	for {

		conn1, err8 := L.Accept()
		if err8 != nil {
			continue
		}
		//если соедниняет провеяет нет ли уже соединения
		if Inchat == false {
			go accepterconn(conn1, txt)
		} else {
			conn1.Close()
		}

	}

}

// дочерний от accepter для непосредсвенного подключения
func accepterconn(conn1 net.Conn, txt *widget.Label) {
	smssr := bufio.NewReader(conn1)
	smss, _ := smssr.ReadString('\n')
	smss = strings.TrimSpace(smss)
	//если получил запрос то спрашивает о подкюченни
	if smss == "12-3-131*-*13" {

		fyne.Do(func() {
			txt.SetText(txt.Text + "\n" + "с вами хотят открыт чат /YM - да /NM - нет")
		})
	}
	////эхо войны
	//if smss == "12-3-131*-*13call" {
	//	println("вам звонят /YMC - да /NMC - нет")
	//	}

	Conn2 = conn1

}

// KEY
func Keyg() {

}

// обрободчик команд когда не в диологе
func master(L net.Listener, gapp fyne.App, txt *widget.Label) {

	fyne.Do(func() {
		txt.SetText(txt.Text + "\n" + "HELP: /CONN - подключится; /ADD - добавить домен; /YM -принять запрос на чат;/EXITM-выйти из чата не закрывая приложение;/CONNCALL")
	})

	var portc string
	var ipc string
	var nameadd string
	var ipadd string
	var portadd string

	go accepter(L, txt)
	for {

		I := <-Keys

		if I == "/HELP" {
			fyne.Do(func() {
				txt.SetText(txt.Text + "\n" + "HELP: /CONN - подключится; /ADD - добавить домен; /YM -принять запрос на чат;/EXITM-выйти из чата не закрывая приложение;/CONNCALL")
			})
			//команада подключения
		} else if I == "/CONN" {
			fyne.Do(func() {
				txt.SetText(txt.Text + "\n" + "напиште /I если хотите подключится по айпи\"")
			})
			I = <-Keys
			if I == "/I" {
				fyne.Do(func() { txt.SetText(txt.Text + "\n" + "ip:") })
				ipc = <-Keys
				fyne.Do(func() { txt.SetText(txt.Text + "\n" + "port:") })
				portc = <-Keys
				var err087 error
				//подключение по айпи
				Conn2, err087 = net.Dial("tcp", ipc+":"+portc)
				if err087 != nil {
					fmt.Println(err087)
					continue
				}
			} else {
				binpart, _ := os.Executable()
				binDir := filepath.Dir(binpart)
				contactdir := filepath.Join(binDir, "Contacts")
				os.MkdirAll(contactdir, 0755)
				fyne.Do(func() { txt.SetText(txt.Text + "\n" + "print name:") })

				nameprint := <-Keys
				contactfile := filepath.Join(contactdir, nameprint+".txt")
				Filename, err := os.OpenFile(contactfile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
				if err != nil {
					println(err.Error())
				}

				ipbuf, err := io.ReadAll(Filename)
				if err != nil {
					println(err.Error())
				}
				Filename.Close()
				sportip := string(ipbuf)
				sportip = strings.TrimSpace(sportip)
				var err087 error
				Conn2, err087 = net.Dial("tcp", sportip)
				if err087 != nil {
					fmt.Println(err087)
					continue
				}

			}

			Conn2.Write([]byte("12-3-131*-*13" + "\n"))
			smssw := bufio.NewReader(Conn2)
			for {
				dogger, err := smssw.ReadString('\n')
				if err != nil {
					println(err.Error())
					continue
				}
				dogger = strings.TrimSpace(dogger)
				if dogger == "12-3-131*-*13+" {
					Inchat = true
					go read(Conn2, txt)
					write(Conn2, txt)
					break
				} else if dogger == "12-3-131*-*13-" {
					Conn2.Close()
					fyne.Do(func() { txt.SetText(txt.Text + "\n" + "отказано") })
					break
				}

			}

		} else if I == "/YM" {
			Conn2.Write([]byte("12-3-131*-*13+" + "\n"))
			Inchat = true
			go read(Conn2, nil)
			write(Conn2, nil)
		} else if I == "/NM" {
			Conn2.Write([]byte("12-3-131*-*13-" + "\n"))
			Conn2.Close()

		} else if I == "/ADD" {
			binpart, _ := os.Executable()
			binDir := filepath.Dir(binpart)
			contactdir := filepath.Join(binDir, "Contacts")

			err := os.MkdirAll(contactdir, 0755)
			if err != nil {
				println(err.Error())
			}
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "name:") })
			nameadd = <-Keys
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "ip:") })
			ipadd = <-Keys
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "port:") })
			portadd = <-Keys
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "name: " + nameadd + " ip: " + ipadd + " port: " + portadd) })
			contactfile := filepath.Join(contactdir, nameadd+".txt")
			Filename, _ := os.OpenFile(contactfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

			scr := strings.NewReader(ipadd + ":" + portadd)
			io.Copy(Filename, scr)
			defer Filename.Close()
			Filename.Close()

		}
		//else if I == "/CONNCALL" {
		//
		//	println("ip:")
		//	I = <-Keys
		//	ipc = I
		//	println("port:")
		//	I = <-Keys
		//	portc = I
		//	var err087 error
		//	Conn2, err087 = net.Dial("tcp", ipc+":"+portc)
		//	if err087 != nil {
		//		fmt.Println(err087)
		//		continue
		//	}
		//	Conn2.Write([]byte("12-3-131*-*13call" + "\n"))
		//	smssw := bufio.NewReader(Conn2)
		//	for {
		//		dogger, err := smssw.ReadString('\n')
		//		if err != nil {
		//			println(err.Error())
		//			break
		//		}
		//		dogger = strings.TrimSpace(dogger)
		//		if dogger == "12-3-131*-*13call+" {
		//			Inchat = true
		//			ports, err33 := smssw.ReadString('\n')
		//			if err33 != nil {
		//				println(err33.Error())
		//			}
		//			ports = strings.TrimSpace(ports)
		//			ipport, err := net.ResolveUDPAddr("udp", ipc+":"+ports)
		//			if err != nil {
		//				println(err.Error())
		//			}
		//
		//			conncallc, _ := net.DialUDP("udp", nil, ipport)
		//			CONNCALLs = true
		//			go readaudio(conncallc)
		//			writeaudio(conncallc, ipport)
		//			break
		//		} else if dogger == "12-3-131*-*13call-" {
		//			Conn2.Close()
		//			println("отказано")
		//			break
		//		}
		//
		//	}
		//} else if I == "/YMC" {
		//	Conn2.Write([]byte("12-3-131*-*13call+" + "\n"))
		//
		//	fmt.Println("updport:")
		//	ports := <-Keys
		//	nru, err := net.ResolveUDPAddr("udp", ":"+ports)
		//	if err != nil {
		//		println(err.Error())
		//	}
		//	lcall, err087 := net.ListenUDP("udp", nru)
		//	if err087 != nil {
		//		fmt.Println(err087)
		//	}
		//	Conn2.Write([]byte(ports + "\n"))
		//
		//	Inchat = true
		//	readaudio(lcall)
		//} else if I == "/NMC" {
		//	Conn2.Write([]byte("12-3-131*-*13call-" + "\n"))
		//	Conn2.Close()
		//}

	}

}

//	func writeaudio(lcall *net.UDPConn, addr *net.UDPAddr) {
//		if CONNCALLs != true {
//			lcall.Write([]byte("/CONNCALLDOG"))
//		}
//
//		inputBuf := make([]int16, frameSize)
//
//		enc, err := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
//		if err != nil {
//			println(err.Error())
//			return
//		}
//		stream, err := portaudio.OpenDefaultStream(channels, 0, sampleRate, frameSize, inputBuf)
//		if err != nil {
//			println(err.Error())
//			return
//		}
//		stream.Start()
//		defer stream.Stop()
//
//		compressed := make([]byte, 1024)
//
//		for {
//			err := stream.Read()
//			if err != nil {
//				continue
//			}
//			n, err := enc.Encode(inputBuf, compressed)
//			if err != nil {
//				continue
//			}
//			if addr != nil {
//				lcall.WriteToUDP(compressed[:n], addr)
//			} else {
//				lcall.Write(compressed[:n])
//			}
//
//		}
//	}
//
//	func readaudio(lcall *net.UDPConn) {
//		hellor := make([]byte, 4096)
//		if CONNCALLs != true {
//			n, remoteAddr, err := lcall.ReadFromUDP(hellor)
//			if err != nil {
//				println(err.Error())
//			}
//			println(hellor[:n])
//			go writeaudio(lcall, remoteAddr)
//		}
//
//		outputBuf := make([]int16, frameSize)
//		dec, err := opus.NewDecoder(sampleRate, channels)
//		if err != nil {
//			println(err.Error())
//
//		}
//		stream, err := portaudio.OpenDefaultStream(0, channels, sampleRate, frameSize, outputBuf)
//		if err != nil {
//			println(err.Error())
//			return
//		}
//		stream.Start()
//		defer stream.Stop()
//		udpduf := make([]byte, 2048)
//
//		for {
//			n, _, err := lcall.ReadFromUDP(udpduf)
//			if err != nil {
//				break
//			}
//			_, err = dec.Decode(udpduf[:n], outputBuf)
//			if err != nil {
//				continue
//			}
//			stream.Write()
//		}
//
// }
// функция писателя
func write(conn net.Conn, txt *widget.Label) {

	for {

		otvet := <-Keys
		if otvet == "/EXITM" {
			Inchat = false
			conn.Close()
			return
		}
		if otvet == "/FILE" {
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "Enter filename") })

			otvet = <-Keys
			var statss, err0 = os.Stat(otvet)
			if err0 != nil {
				fmt.Println("file not exist" + err0.Error())
				continue
			}
			conn.Write([]byte("/FILE" + "\n"))
			conn.Write([]byte(fmt.Sprintf("%d", statss.Size()) + "\n"))
			conn.Write([]byte(statss.Name() + "\n"))
			var file, err1 = os.Open(otvet)
			if err1 != nil {
				fmt.Println("file not exist" + err1.Error())
			}
			io.Copy(conn, file)
			file.Close()

		} else {
			conn.Write([]byte(otvet + "\n"))
		}

	}

}

// функция читателя
func read(conn net.Conn, txt *widget.Label) {

	bufer := bufio.NewReader(conn)

	for {
		vopros, err3 := bufer.ReadString('\n')
		if err3 != nil {
			println(err3.Error())
			continue
		}
		vopros = strings.TrimSpace(vopros)
		if vopros == strings.TrimSpace("/FILE") {

			sizef, err2 := bufer.ReadString('\n')
			if err2 != nil {
				println(err2.Error())
				continue
			}
			var namef, err4 = bufer.ReadString('\n')
			if err4 != nil {
				println(err4.Error())
				continue
			}

			num, err5 := strconv.ParseFloat(strings.TrimSpace(sizef), 64)
			if err5 != nil {
				println(err5.Error())
				continue
			}

			num = num / 1048576
			snum := strconv.FormatFloat(num, 'f', -1, 64)

			println("download file name- " + strings.TrimSpace(namef) + " file size- " + snum + " Mbyte")
			fyne.Do(func() {
				txt.SetText(txt.Text + "\n" + "download file name- " + strings.TrimSpace(namef) + " file size- " + snum + " Mbyte")
			})

			binpart, _ := os.Executable()
			binDir := filepath.Dir(binpart)
			downloaddir := filepath.Join(binDir, "DownLoad")
			err := os.MkdirAll(downloaddir, 0755)
			if err != nil {
				println(err.Error())
				continue
			}

			f, err6 := os.Create(filepath.Join(downloaddir, strings.TrimSpace(namef)))
			if err6 != nil {
				println(err6.Error())
				continue
			}

			scr, err7 := strconv.ParseInt(strings.TrimSpace(sizef), 10, 64)
			if err7 != nil {
				println(err7.Error())
				continue
			}
			io.CopyN(f, bufer, scr)
			f.Close()
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "file- \" + strings.TrimSpace(namef) + \" downloaded") })
		} else {
			fyne.Do(func() { txt.SetText(txt.Text + "\n" + "file- \" + strings.TrimSpace(namef) + \" downloaded") })
		}

	}

}

// ининициализатор
func main() {
	//err := portaudio.Initialize()
	//if err != nil {
	//	fmt.Printf("init portaudio err: %v\n", err)
	//}
	//defer portaudio.Terminate()
	//file, err2 := os.Open("contacts.txt")
	//if err2 != nil {
	//	filel, err := os.Create("contacts.txt")
	//	if err != nil {
	//		println(err.Error())
	//		return
	//	}
	//	filel.Close()
	//}
	//file.Close()
	//читает и открывает слушателя
	//var ports string
	//println("write you port")
	//fmt.Scanln(&ports)
	L, err := net.Listen("tcp", ":25656")
	if err != nil {
		println(err.Error())
	}

	//начало выполнения функции меню(не в диологе)

	App := app.New()
	awin := App.NewWindow("Phillichat")
	var text string
	txt := widget.NewLabel(text)
	ent := widget.NewEntry()
	sendbut := widget.NewButton("sendbut", func() {
		Keys <- ent.Text
		txt.SetText(txt.Text + "\n" + ent.Text)
		ent.Text = ""
	})
	awin.Resize(fyne.NewSize(300, 1200))
	awin.SetFullScreen(true)
	awin.SetContent(container.NewVBox(txt, ent, sendbut))
	//запукает основную программу
	go master(L, App, txt)
	go Keyg()
	awin.ShowAndRun()
}
