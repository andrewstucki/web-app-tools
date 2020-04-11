import React from "react";
import { shallow } from "enzyme";
import toJson from "enzyme-to-json";
import { User } from "./models";

const injectState = (mockUser: User | null) => {
  const mockDispatch = jest.fn();
  jest.mock("react-redux", () => ({
    ...jest.requireActual("react-redux"),
    useSelector: (): User | null => mockUser,
    useDispatch: () => mockDispatch,
  }));
  const App = require("./App").App;
  return { App, dispatch: mockDispatch };
};

describe("App", () => {
  beforeEach(() => {
    jest.resetModules();
  });

  it("renders loading state with no user", () => {
    const { App } = injectState(null);
    expect(toJson(shallow(<App />))).toMatchSnapshot();
  });

  it("requests the user's profile", () => {
    const { App, dispatch } = injectState(null);
    const node = toJson(shallow(<App />));
    expect(dispatch).toBeCalledWith(
      expect.objectContaining({
        type: "__AUTHENTICATED_REQUEST",
        payload: expect.objectContaining({
          config: { url: "/api/v1/me" },
        }),
      })
    );
  });

  it("renders basic info given a user", () => {
    const { App } = injectState({
      id: "foo",
      email: "bar@baz.com",
      createdAt: new Date(),
      updatedAt: new Date(),
    });
    expect(toJson(shallow(<App />))).toMatchSnapshot();
  });

  it("logs out properly", () => {
    const { App, dispatch } = injectState({
      id: "foo",
      email: "bar@baz.com",
      createdAt: new Date(),
      updatedAt: new Date(),
    });
    const node = shallow(<App />)
      .find("button")
      .simulate("click");
    expect(dispatch).toBeCalledWith(
      expect.objectContaining({
        type: "__AUTHENTICATED_LOG_OUT",
      })
    );
  });
});
