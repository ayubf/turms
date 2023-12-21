import {
    createBrowserRouter,
  } from "react-router-dom";
import Home from "./pages/Home";
import Room from "./pages/Room";
import { RenderErrorBoundary } from "react-router/dist/lib/hooks";

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
    }
  ]);
  

  export default TurmsRouter