import React, { Component } from 'react';

export default class ChildHello extends Component {
  render() {
    const { message = 'no message' } = this.props;
    return (
      <span className="something">
        <strong>
          message: {message}
        </strong>
      </span>
    );
  }
}
