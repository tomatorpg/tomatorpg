import React, { Component } from 'react';
import ChildHello from './ChildHello';

export default class Hello extends Component {
  render() {
    return (
      <div>
        <ChildHello message="Hello World" />
      </div>
    );
  }
}
