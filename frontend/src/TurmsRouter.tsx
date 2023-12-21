import {
    createBrowserRouter,
  } from "react-router-dom";
import Home from "./pages/Home";
import Room from "./pages/Room";
import CreateUserName from "./pages/CreateUserName";

const TurmsRouter = createBrowserRouter([
    {
      path: "/",
      element: <Home />
    },
    {
      path: "/joinroom/:id",
      element: <Room />,
      errorElement: (<div>
        <h1>404: Room Not Found</h1>
        <p>Room may have expired or room code was inputted incorrectly</p>
      </div>)
    },
    {
      path: "/createuser/:id",
      element: <CreateUserName />,
      errorElement: (<div>
        <h1>404: Room Not Found</h1>
        <p>Room may have expired or room code was inputted incorrectly</p>
      </div>)
    }, 
    {
      path: "/roomclosed", 
      element: (<div>
        <h1>404: Room Not Found</h1>
        <p>Room may have expired or room code was inputted incorrectly</p>
      </div>)
    }
  ]);
  

  export default TurmsRouter