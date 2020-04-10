import React from "react";
import { Provider } from 'react-redux';
import { mount } from 'enzyme';
import { createStore } from 'redux';
import { reducer, initialState } from './store';
import { App } from './App';

describe('App', () => {
  const mockStore = createStore(reducer, initialState);
  mockStore.dispatch = jest.fn();

  it('Renders without crashing', () => {
    const app = mount(
      <Provider store={mockStore}>
        <App/>
      </Provider>
    );
    expect(app).toMatchSnapshot();
  });
});
