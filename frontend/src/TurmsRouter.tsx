import {
    createBrowserRouter,
  } from "react-router-dom";
import Home from "./pages/Home";
import Room from "./pages/Room";

const TurmsRouter = createBrowserRouter([
    {
      path: "/",
      element: <Home />
    },
    {
      path: "/joinroom",
      element: <Room />
    }
  ]);
  

  export default TurmsRouter