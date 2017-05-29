import React from 'react';
import { Link } from 'react-router-dom';

const Nav = (state) => {
  const { user } = state;
  const userLink = (user.id === 0) ? (
    <a href="/oauth2/google">Login</a>
  ) : (
    <a href="/oauth2/logout">Logout</a>
  );
  return (
    <nav className="main-menu">
      <ul>
        <li className="user">Hello, <span className="name">{ user.name }</span></li>
        <li><Link to="/rooms/">Rooms</Link></li>
        <li>{ userLink }</li>
      </ul>
    </nav>
  );
};

export default Nav;
