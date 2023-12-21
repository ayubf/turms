import React, {useEffect, useState} from 'react';
import { useParams } from 'react-router';
import { useNavigate } from 'react-router-dom';



interface Message {
    Author: string,
    Message: string,
    Time?: String, 
}

interface RoomData {
    roomCreatorName: string, 
    roomName: string, 
    code: string,
    expirationAt: string,
    messages: Array<Message>
}

const Room = () => {
    const nav = useNavigate();
    const {id} = useParams();
    const [sessionExists, setSessionExist] = useState(false);
    const [username, setUsername] = useState("");
    const [messageList, setMessageList] = useState(Array<Message>());
    const [roomData, setRoomData] = useState<RoomData>();
    const [currMessage, setCurrMessage] = useState("");
    const [socket, setSocket] = useState<WebSocket>();



    const checkSession = async () => {
        await fetch("http://localhost:8080/getsession", {
            credentials: 'include' 
        })
        .then(i => i.json())
        .then((data) => {
            if (data.hasOwnProperty("iss")) {
                setSessionExist(true); 
                setUsername(data["iss"])
                console.log("WOW")
            } else {
                nav(`/createuser/${id}`)
            }
        })
    }

    const websocketConnect = async () => {
        const socket = new WebSocket(`ws://localhost:8080/joinroom?code=${id}`);
        console.log("Trying to connect.....");

        socket.onopen = () => {
            console.log("Connection successful!");

            socket.onmessage = async (e) => {
                console.log("Received message from server:", e.data);
                if (!JSON.parse(e.data).hasOwnProperty("Message")) {
                    const r = JSON.parse(e.data);
                    setRoomData(r)
                    setMessageList(r.messages)
                }
            };
        };

        socket.onerror = () => {
            nav("/roomclosed")
        }

        socket.onclose = () => {
            console.log("Closing connection....");
            setMessageList([...messageList, {"Message": "Room Closed", "Author": "Turms"}])
        };

        setSocket(socket);

    };

    const sendMessage = (message: string) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            const messageJSON = JSON.stringify({ "Message": message });
            socket.send(messageJSON);
            setMessageList([...messageList, {"Message": message, "Author": username}])
        }
    };


    useEffect(() => {
        const fetchData = async () => {
            await checkSession(); 
    
            if (sessionExists) {
                websocketConnect(); 
            }
        };
    
        fetchData();
        return () => {}
    }, [sessionExists, id]);


    return <div>
        {
            roomData != undefined ? (
                <div>
                    <h1>{roomData.roomName}</h1>
                    <h3>{roomData.roomCreatorName}</h3>
                    <p>Expires at: {roomData.expirationAt}</p>
                    <button
                        onClick={() => {
                            navigator.clipboard.writeText(`http://localhost:3000/joinroom/${roomData.code}`);
                        }}
                    >Copy Room Link</button>
                </div>
            ) : <></>
        }
        <div className="messageArea">
            {
                roomData ? (
                    <div>
                        <div>
                            {
                                messageList.map( ({Author, Message, Time}, i) => {
                                    return <p style={{
                                        "color": Author == "Turms" ? "red" : "black"
                                    }} key={i} >{Author} : {Message}</p>
                                })
                            }
                        </div>
                        <input type="text" value={currMessage} onChange={(e) => setCurrMessage(e.target.value)} />
                        <button onClick={() => {sendMessage(currMessage)}}>Send Message</button>
                    </div>
                ) : <></>
            } 
        </div>
    </div>
}

export default Room;