/* eslint
react/forbid-prop-types: 'warn'
*/
import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Rooms from './Rooms';
import Room from './Room';

class App extends Component {
  render() {
    const { dispatch, roomActivities, rooms } = this.props;
    return (
      <div>
        <Rooms dispatch={dispatch} rooms={rooms} />
        <Room dispatch={dispatch} roomActivities={roomActivities} />
      </div>
    );
  }
}

App.propTypes = {
  dispatch: PropTypes.func,
  rooms: PropTypes.array,
  roomActivities: PropTypes.array,
};

App.defaultProps = {
  dispatch: () => {},
  rooms: [],
  roomActivities: [],
};

export default App;
