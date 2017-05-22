/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { messageInRoom } from '../transports/JSONSocket';

class Room extends Component {

  componentDidUpdate() {
     // scroll to bottom after any updates
    this.messageWrapper.scrollTop = this.messageBox.clientHeight;
  }

  submitHandler(evt) {
    const { dispatch } = this.props;

    evt.preventDefault(); // prevent form submission

    // send message with server object
    dispatch(messageInRoom(this.textInput.value));

    // reset text box
    this.textInput.value = '';
  }

  render() {
    const { roomActivities } = this.props;
    const messagesSummary = (roomActivities.length > 0) ?
      roomActivities.reduce((acc, activity, index) => {
        const { type, message = '' } = activity;
        if (type === 'message') {
          const key = `message-${index}`;
          acc.push(<div key={key}>{message}</div>);
        }
        return acc;
      }, []) :
      <div>No message yet</div>;
    return (
      <div id="room">
        <div className="messages-wrapper" ref={(element) => { this.messageWrapper = element; }}>
          <div className="messages" ref={(messages) => { this.messageBox = messages; }}>
            {messagesSummary}
          </div>
        </div>
        <form onSubmit={evt => this.submitHandler(evt)}>
          <input type="text" ref={(input) => { this.textInput = input; }} />
          <button type="submit">Submit</button>
        </form>
      </div>
    );
  }
}

Room.propTypes = {
  dispatch: PropTypes.func,
  roomActivities: PropTypes.array,
};

Room.defaultProps = {
  dispatch: () => {},
  roomActivities: [],
};

export default Room;
