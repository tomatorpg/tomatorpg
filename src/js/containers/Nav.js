import React, { Component } from 'react';

class Nav extends Component {
  render() {
    return (
      <nav className="main-menu">
        <ul>
          <li><a href="/oauth2/google">Login</a></li>
          <li><a href="/oauth2/logout">Logout</a></li>
        </ul>
      </nav>
    );
  }
}

export default Nav;
