import {useState} from "react";
import { useNavigate } from "react-router";

const Join = () => {

    const [code, setCode] = useState("");
    const nav = useNavigate();

    const joinRoom = () => {
        nav(`/joinroom/${code}`)
    }

    return (
        <div>
            <h1>Join Room</h1>
            <hr />
            <p>Enter a room code to join a new room or <a href="/">create one</a></p>
            <div className="fieldWrapper">
                <input type="text" value={code} onChange={(e) => { setCode(e.target.value) }} id="username" />
            </div>
            <div className="fieldWrapper">
          <button onClick={() => { joinRoom() }}>Join</button>
        </div>
        </div>
    )

}

export default Join;