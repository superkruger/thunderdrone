import {NavLink} from "react-router-dom";

export default function Navbar(props) {

    return <>
        <NavLink to="/" className="nav-brand">
            Thunderdrone
        </NavLink>
        <ul>
            <li className="nav-item">
                <NavLink exact activeClassName="active" to="/">Dashboard</NavLink>
            </li>
            <li className="nav-item">
                <NavLink exact activeClassName="active" to="/settings">Settings</NavLink>
            </li>
        </ul>
    </>
}