import { Outlet } from "react-router-dom";
import "./styles.scss";

export default function Root() {

  return (
    <>
      <div className="middle">
        <Outlet/>
      </div>
    </>
  );
}
