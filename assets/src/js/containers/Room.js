/* eslint
react/forbid-prop-types: 'warn',
jsx-a11y/no-autofocus: 'warn',
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { messageInRoom } from '../transports/JSONSocket';

class Room extends Component {

  componentDidMount() {
    const { onLoad } = this.props;
    onLoad(this.props);
  }

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
    const { roomActivities, listCharacters } = this.props;
    console.log('listCharacters', listCharacters);
    const messagesSummary = (roomActivities.length > 0) ?
      roomActivities.reduce((acc, activity, index) => {
        const { type, message: { message = '', userID = 0 } } = activity;
        if (type === 'message') {
          const key = `message-${index}`;
          const userDisplayName = (userID === 0) ? 'Visitor' : `User ${userID}`;
          acc.push(<div className="message-wrapper" key={key}>
            <div className="user">
              {userDisplayName}
            </div>
            <div className="message">
              {message}
            </div>
          </div>);
        }
        return acc;
      }, []) :
      <div className="no-message">No message yet</div>;
    return (
      <div className="container">
        <div id="room">
          <div className="messages-wrapper" ref={(element) => { this.messageWrapper = element; }}>
            <div className="area-header">
              Chat
            </div>
            <div className="messages" ref={(messages) => { this.messageBox = messages; }}>
              {messagesSummary}
            </div>
          </div>
          <form className="room-form" onSubmit={evt => this.submitHandler(evt)}>
            <div className="textarea-wrapper">
              <div className="area-header">
                Your Messages
              </div>
              <textarea type="text" autoFocus ref={(input) => { this.textInput = input; }} />
              <div className="actions">
                <button className="btn" type="button" onClick={() => listCharacters()}>Character</button>
                <button className="btn btn-primary" type="submit">Send</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    );
  }
}

Room.propTypes = {
  dispatch: PropTypes.func,
  listCharacters: PropTypes.func,
  onLoad: PropTypes.func,
  roomActivities: PropTypes.array,
};

Room.defaultProps = {
  dispatch: () => {},
  listCharacters: () => {},
  onLoad: () => {},
  roomActivities: [],
};

export default Room;
