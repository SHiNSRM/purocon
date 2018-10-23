package main

import (
    "fmt"
    "net/http"
    "log"
    "math/rand"
    "time"
    "strconv"
    "strings"
)
type String string
// http.HandleFuncに登録する関数
// http.ResponseWriterとhttp.Requestを受ける
var user=make([][]int,12)
var pcalc=make([][]int,12)
var field=make([][]int,12)
var turn=0
var length=0
var width=0
var p=make(map[int]map[string]int)
var pcount [5]int = [5]int{0, 0, 0, 0, 0}

func StartServer(w http.ResponseWriter, r *http.Request) {
    rand.Seed(time.Now().UnixNano())
    turn=rand.Intn(60)+60
    length=rand.Intn(4)+8
    width=rand.Intn(4)+8
    fmt.Fprintf(w,"%d\n",turn)
    fmt.Fprintf(w,"%d\n",length)
    fmt.Fprintf(w,"%d\n",width)
    field=make([][]int,(length+1)/2)
    for i:=0; i<(length+1)/2; i++{
      field[i]=make([]int, width)
      for j:=0; j<width; j++ {
        field[i][j]=rand.Intn(32)-16
        //fmt.Fprintf(w,"%d ",field[i][j])
      }
      //fmt.Fprintf(w,"\n")
    }

    tmp_field:=make([][]int,length/2)
    for i:=0; i<length/2; i++{
      tmp_field[i]=make([]int, width)
      tmp_field[i]=field[((length)/2)-1-i]
    }
    field=append(field,tmp_field...)

    for i:=0; i<length; i++{
      for j:=0; j<width; j++ {
        fmt.Fprintf(w,"%d ",field[i][j])
      }
      fmt.Fprintf(w,"\n")
    }

    //user:=make([][]int,length)
    for i:=0; i<length; i++{
      user[i]=make([]int, width)
    }
/*
    fmt.Fprintf(w,"%d ",width)
    fmt.Fprintf(w,"\n")
    fmt.Fprintf(w,"%d ",width/2-1)
    fmt.Fprintf(w,"\n")
    fmt.Fprintf(w,"%d ",(width/2-1)-2)
    //a:=rand.Intn((width/2-1)-2)+1
    //fmt.Fprintf(w,"%d ",a)
    fmt.Fprintf(w,"\n")
*/
    //p:=make(map[int]map[string]int)
    for i:=1; i<5; i++{
      p[i]=make(map[string]int)
    }
    x:=rand.Intn((width/2-1)-2)+1
    y:=rand.Intn((length/2-1)-2)+1
    p[1]["x"]=x
    p[1]["y"]=y
    p[2]["x"]=x
    p[2]["y"]=width-y-1
    p[3]["x"]=length-x-1
    p[3]["y"]=y
    p[4]["x"]=length-x-1
    p[4]["y"]=width-y-1

    for i:=1; i<5; i++{
      user[p[i]["x"]][p[i]["y"]]=i
    }

/*
    user[x][y]=1
    user[x][width-y-1]=2
    user[length-x-1][y]=3
    user[length-x-1][width-y-1]=4
*/
    for i:=0; i<length; i++{
      for j:=0; j<width; j++ {
        fmt.Fprintf(w,"%d ",user[i][j])
      }
      fmt.Fprintf(w,"\n")
    }

}

func MoveServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "move\n")
    r.ParseForm()
    //curl -X POST localhost:8000/move -d "usr=1&d=right"
    u,_:=strconv.Atoi(r.FormValue("usr"))
    fmt.Println(u)
    fmt.Println(r.FormValue("d"))
    //d:=r.FormValue("d")
    d:=strings.Split(r.FormValue("d"), "")
    fmt.Println(d);
    /*
    for i:=0; i<len(d); i++{
      if d[i]=="r"{p[u]["y"]++
      }else if d[i]=="l"{p[u]["y"]--
      }else if d[i]=="u"{p[u]["x"]--
      }else if d[i]=="d"{p[u]["x"]++}
    }
    */
    tmp_px:=p[u]["x"]
    tmp_py:=p[u]["y"]
    for i:=0; i<len(d); i++{
      if d[i]=="r"{tmp_py++
      }else if d[i]=="l"{tmp_py--
      }else if d[i]=="u"{tmp_px--
      }else if d[i]=="d"{tmp_px++}
    }
    if 0<=tmp_px && tmp_px<length && 0<=tmp_py && tmp_py<width {
      if u==1||u==2 {
        if user[tmp_px][tmp_py]==0 || user[tmp_px][tmp_py]==5 {
          user[p[u]["x"]][p[u]["y"]]=5
        }else{
          fmt.Fprintf(w,"Error \n")
          return
        }
      }else{
        if user[tmp_px][tmp_py]==0 || user[tmp_px][tmp_py]==6 {
          user[p[u]["x"]][p[u]["y"]]=6
        }else{
          fmt.Fprintf(w,"Error \n")
          return
        }
      }
      p[u]["x"]=tmp_px
      p[u]["y"]=tmp_py
    }else{
      fmt.Fprintf(w,"Error \n")
      return
    }
    user[p[u]["x"]][p[u]["y"]]=u
    pcount[u]++
    if(pcount[1]==pcount[2]&&pcount[2]==pcount[3]&&pcount[3]==pcount[4]){
      pcount[0]=pcount[1]
      fmt.Fprintf(w,"%d ",pcount[0])
    }
    if(turn==pcount[0]){
      fmt.Fprintf(w,"end the game \n")
    }
}

func RemoveServer(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "remove\n")
  r.ParseForm()
  //curl -X POST localhost:8000/move -d "usr=1&d=right"
  u,_:=strconv.Atoi(r.FormValue("usr"))
  fmt.Println(u)
  fmt.Println(r.FormValue("d"))
  d:=strings.Split(r.FormValue("d"), "")
  tmp_px:=p[u]["x"]
  tmp_py:=p[u]["y"]
  for i:=0; i<len(d); i++{
    if d[i]=="r"{tmp_py++
    }else if d[i]=="l"{tmp_py--
    }else if d[i]=="u"{tmp_px--
    }else if d[i]=="d"{tmp_px++}
  }
  if 0<=tmp_px && tmp_px<length && 0<=tmp_py && tmp_py<width {
    if user[tmp_px][tmp_py]!=1&&user[tmp_px][tmp_py]!=2&&user[tmp_px][tmp_py]!=3&&user[tmp_px][tmp_py]!=4 {user[tmp_px][tmp_py]=0}
  }else{
    fmt.Fprintf(w,"Error \n")
    return
  }

  pcount[u]++
  if(pcount[1]==pcount[2]&&pcount[2]==pcount[3]&&pcount[3]==pcount[4]){
    pcount[0]=pcount[1]
    fmt.Fprintf(w,"%d ",pcount[0])
  }
  if(turn==pcount[0]){
    fmt.Fprintf(w,"end the game \n")
  }
}

func ShowServer(w http.ResponseWriter, r *http.Request) {
  for i:=0; i<length; i++{
    for j:=0; j<width; j++ {
      fmt.Fprintf(w,"%d ",field[i][j])
    }
    fmt.Fprintf(w,"\n")
  }
  for i:=0; i<length; i++{
    for j:=0; j<width; j++ {
      fmt.Fprintf(w,"%d ",user[i][j])
    }
    fmt.Fprintf(w,"\n")
  }
}

func UsrpointServer(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "usrpoint\n")
  r.ParseForm()
  u,_:=strconv.Atoi(r.FormValue("usr"))
  fmt.Println(p[u]["x"])
  fmt.Println(p[u]["y"])
  fmt.Fprintf(w,"%d ",p[u]["y"])
  fmt.Fprintf(w,"%d",p[u]["x"])
}

func myAbs(x int) int{
  if(x<0){return -x}
  return x
}

var use5[60][60] bool
var use6[60][60] bool
var came[60][60] bool
var dx [4]int = [4]int{1, 0, -1, 0}
var dy [4]int = [4]int{0, 1, 0, -1}
var flag bool
var cnt int
func check_area(y int,x int ,wall int)bool{

  ret:=true
  came[y][x]=true
  if(!flag){return false}
  if(pcalc[y][x]==wall){return true}
  for i:=0;i<4;i++{
    nx:=x+dx[i]
    ny:=y+dy[i]
    tmp:=true
    if(nx<0||ny<0||nx>=width||ny>=length){
      flag=false
      return false
    }
    if(!came[ny][nx]){tmp=check_area(ny,nx,wall)}
    if(!tmp){ret=false}
  }
  return ret
}

func init_check_area(){
  flag=true
  for i:=0;i<length;i++{
    for j:=0;j<width;j++{
      came[i][j]=false
    }
  }
}

func PointcalcServer(w http.ResponseWriter, r *http.Request) {
  pcalc=user
  point5:=0
  point6:=0


  for i:=0; i<length; i++{
    for j:=0; j<width; j++ {
      if(pcalc[i][j]==1||pcalc[i][j]==2){
        pcalc[i][j]=5
      }
      if(pcalc[i][j]==3||pcalc[i][j]==4){
        pcalc[i][j]=6
      }
      fmt.Fprintf(w,"%d ",pcalc[i][j])
    }
    fmt.Fprintf(w,"\n")
  }

  fmt.Fprintf(w,"盤面\n")
  for i:=0; i<length; i++{
    for j:=0; j<width; j++ {
    fmt.Fprintf(w,"%04d ",field[i][j])
    }
    fmt.Fprintf(w,"\n")
  }

  //////////以上プリントでバッグ
  for i:=0;i<length;i++{
    for j:=0;j<width;j++{
      use5[i][j]=false
      use6[i][j]=false
    }
  }
  for i:=0;i<length;i++{
    for j:=0;j<width;j++{
      init_check_area() //flag=trueにして、cameをすべてfalseにする
      if(check_area(i,j,5)&&!use5[i][j]){
        use5[i][j]=true;

      }
      init_check_area()
      if(check_area(i,j,6)&&!use6[i][j]){
        use6[i][j]=true;

      }
    }
  }

  for y:=0;y<length;y++{//縦
    for x:=0;x<width;x++{//横
      if(use5[y][x]){
        if(pcalc[y][x]==5){point5+=field[y][x]
        }else{point5+=myAbs(field[y][x])}
      }
      if(use6[y][x]){
        if(pcalc[y][x]==6){point6+=field[y][x]
        }else{point6+=myAbs(field[y][x])}
      }

    }
  }
  fmt.Fprintf(w,"score = %d ",point5)
  fmt.Fprintf(w,"score = %d ",point6)

}


func InitServer(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  //u,_:=strconv.Atoi(r.FormValue())
  fmt.Fprintf(w,"%d ",r.ParseForm())
}

func main() {
    // http.HandleFuncにルーティングと処理する関数を登録
    http.HandleFunc("/start", StartServer)
    http.HandleFunc("/move", MoveServer)
    http.HandleFunc("/remove", RemoveServer)
    http.HandleFunc("/show", ShowServer)
    http.HandleFunc("/usrpoint", UsrpointServer)
    http.HandleFunc("/pointcalc", PointcalcServer)
    http.HandleFunc("/init", InitServer)

    // ログ出力
    log.Printf("Start Go HTTP Server")

    // http.ListenAndServeで待ち受けるportを指定
    err := http.ListenAndServe(":8000", nil)

    // エラー処理
    if err != nil {
       log.Fatal("ListenAndServe: ", err)
    }
}
