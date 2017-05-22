/* eslint
react/forbid-prop-types: 'warn'
*/
import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Rooms from './Rooms';
import Room from './Room';

class App extends Component {
  render() {
    const { dispatch, roomActivities } = this.props;
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
        <Rooms dispatch={dispatch} rooms={rooms} />
        <Room dispatch={dispatch} roomActivities={roomActivities} />
      </div>
    );
  }
}

App.propTypes = {
  dispatch: PropTypes.func,
  roomActivities: PropTypes.array,
};

App.defaultProps = {
  dispatch: () => {},
  roomActivities: [],
};

export default App;
