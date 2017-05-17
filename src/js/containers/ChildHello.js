import React, { Component } from 'react';
import PropTypes from 'prop-types';

export default class ChildHello extends Component {

  static propTypes = {
    message: PropTypes.string,
  }

  static defaultProps = {
    message: 'no message',
  }

  render() {
    const { message } = this.props;
    return (
      <span className="something">
        <strong>
          message: {message}
        </strong>
      </span>
    );
  }
}
