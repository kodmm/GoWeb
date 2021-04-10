package main 

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	//leaveはチャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	//clientsには在室している全てのクライアントが保持されます。
	clients map[*client]bool
}

func (r *room) run() {
	for{
		select {
		case client := <- r.join:
			// 参加
			r.clients[client] = true
		case client := <- r.leave:
			// 不参加
			delete(r.clients, client)
			close(client.send)
		case msg := <- r.forward:
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
