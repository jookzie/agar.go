import { createBrowserRouter } from "react-router-dom";
import ErrorPage from "../pages/ErrorPage.tsx";
import Root from "../pages/Root.tsx";
import MainPage from "../pages/MainPage/MainPage.tsx";

const Router = createBrowserRouter([
    {
        path: "/",
        element: <Root/>,
        errorElement: <ErrorPage/>,
        children: [
            {
                path: "/",
                element: <MainPage/>,
            },
        ]
    }
]);
export default Router;
