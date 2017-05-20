/* eslint
react/forbid-prop-types: 'warn'
*/
import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Rooms from './Rooms';
import Room from './Room';

class App extends Component {
  render() {
    const { server, roomActivities } = this.props;
    console.log('app render: this.props', this.props);
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
        <Room server={server} roomActivities={roomActivities} />
      </div>
    );
  }
}

App.propTypes = {
  server: PropTypes.object,
  roomActivities: PropTypes.array,
};

App.defaultProps = {
  server: {},
  roomActivities: [],
};

export default App;
