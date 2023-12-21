import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';


function Home() {

    const navigate = useNavigate();

    function getRandomInt() {
        return Math.floor(Math.random() * 99999);
    }
    
    const [sessionExists, setSessionExist] = useState(false);
    const [username, setUsername] = useState("");
    const [fetchedUser, setFetchedUser] = useState("");
    const [roomname, setRoomname] = useState("Room"+getRandomInt())
    const [timeLimit, setTimeLimit] = useState(15);

    const checkSession = async () => {
        await fetch(`http://localhost:8080/getsession`, {
            credentials: 'include' 
        })
        .then(i => i.json())
        .then((data) => {
            if (data.hasOwnProperty("iss")) {
                setSessionExist(true); 
                setUsername(data["iss"])
                setFetchedUser(data["iss"])
            } else {
                setUsername("Guest"+getRandomInt())
            }
        })
    }

    useEffect(() => {
        checkSession()
        console.log(sessionExists)
    })

    async function createRoom() {
        let res; 

        if (!sessionExists || fetchedUser != username) {
            await fetch(`http://localhost:8080/createsession`, {
                method: "POST",
                body: JSON.stringify({
                    "username": username
                }),
                credentials: 'include' 
            })
        }

        await fetch("http://localhost:8080/createroom", {
            method: "POST",
            body: JSON.stringify({
                "roomname": roomname,
                "timeLimit": timeLimit
            }),
            credentials: 'include' 
        })
        .then(i => i.json())
        .then(data => {
            res = data["code"]
        })
        navigate(`/joinroom/${res}`);
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
                <select name="roomTimer" id="roomTimer" value={timeLimit} onChange={(e) => {setTimeLimit(parseInt(e.target.value))}}>
                    <option value="15">15 Minutes</option>
                    <option value="30">30 Minutes</option>
                    <option value="60">60 Minutes</option>
                    <option value="90">90 Minutes</option>
                </select>
            </div>
            <div className="fieldWrapper">
                <button onClick={() => {createRoom()}}>Create Room</button>
            </div>
        </div>
    </div>
}

export default Home