import { useState } from "react"

function Home() {

    function getRandomInt() {
        return Math.floor(Math.random() * 99999);
    }
    
    const [username, setUsername] = useState("Guest"+getRandomInt());
    const [roomname, setRoomname] = useState("Room"+getRandomInt())

    async function createRoom() {
        let res; 
        await fetch("http://localhost:8080/createroom", {
            method: "POST",
            body: JSON.stringify({
                "username":  username,
                "roomname": roomname
            })
        })
        .then(i => i.json())
        .then(data => res = data)

        console.log(res)
    }

    return <div>
        <div className="createBox">
            <div className="fieldWrapper">
                <input type="text" value={username} onChange={(e) => {setUsername(e.target.value)}} id="username" />
            </div>
            <div className="fieldWrapper">
                <input type="text" value={roomname} onChange={(e) => {setRoomname(e.target.value)}} id="roomname" />
            </div>
            <div className="fieldWrapper">
                <select name="roomTimer" id="roomTimer">
                    <option value={15}>15 Minutes</option>
                    <option value={30}>30 Minutes</option>
                    <option value={60}>60 Minutes</option>
                    <option value={90}>90 Minutes</option>
                </select>
            </div>
            <div className="fieldWrapper">
                <button onClick={() => {createRoom()}}>Create Room</button>
            </div>
        </div>
    </div>
}

export default Home