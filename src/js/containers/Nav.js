import React from 'react';
import { Link } from 'react-router-dom';

const Nav = () => (
  <nav className="main-menu">
    <ul>
      <li><Link to="/rooms/">Rooms</Link></li>
      <li><a href="/oauth2/google">Login</a></li>
      <li><a href="/oauth2/logout">Logout</a></li>
    </ul>
  </nav>
);

export default Nav;
