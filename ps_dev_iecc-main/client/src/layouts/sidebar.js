import { useNavigate } from "react-router-dom";
import Logo from "../assets/img/logo.png";
import { appBase } from "../utils/settings";

function AppSidebar(props) {
  const navigate = useNavigate();
  return (
    <div
      className={`${
        props.open ? "fixed left-0 z-20 shadow-lg" : "hidden"
      } w-max w-25 min-w-25 bg-secondary transition-left duration-300 overflow-auto flex flex-col items-center p-5 md:flex md:relative`}
    >
      <img width={40} src={Logo} alt="logo" />

      <div
        className={`${
          props.open ? "pt-5" : ""
        } flex-1 flex flex-col items-center justify-center gap-8`}
      >
        {props.menu.map(
          (menu, i) =>
            menu.menu && (
              <div
                key={i}
                onClick={() => {
                  if (menu.path.includes("https://")) {
                    window.location.href = menu.path;
                  } else {
                    navigate(appBase + menu.path);
                    props.onSelectMenu(menu.id);
                    props.handleSideBar();
                  }
                }}
                className={`group flex gap-2 items-center p-2 pb-1 rounded-lg cursor-pointer hover:bg-primary hover:text-white ${
                  props.selectedmenu === menu.id
                    ? "bg-primary text-white"
                    : "text-iconColor"
                }`}
              >
                <i className={"bx " + menu.icon} style={{ fontSize: 28 }}></i>
                <h2
                  className={`hidden z-50	fixed ml-[25px] mb-1 p-2 pl-3 pt-2 rounded-lg bg-primary group-hover:block`}
                >
                  {menu.name}
                </h2>
              </div>
            )
        )}
      </div>
    </div>
  );
}

export default AppSidebar;
