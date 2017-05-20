/* eslint
react/forbid-prop-types: 'warn'
*/
import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Rooms from './Rooms';
import Room from './Room';

class App extends Component {
  render() {
    const { server } = this.props;
    const rooms = [
      {
        name: 'hello room 1',
      },
      {
        name: 'hello room 1',
      },
    ];
    return (
      <div>
        <Rooms rooms={rooms} />
        <Room server={server} />
      </div>
    );
  }
}

App.propTypes = {
  server: PropTypes.object,
};

App.defaultProps = {
  server: {},
};

export default App;
