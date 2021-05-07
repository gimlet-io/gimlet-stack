import React, {Component} from 'react'
import {hot} from 'react-hot-loader'
import * as stackDefinitionFixture from '../fixtures/stack-definition.json'

import HelmUI from 'helm-react-ui'
import './style.css'
import StreamingBackend from './streamingBackend'
import GimletCLIClient from './client'
import {Tile} from "./tile";

class App extends Component {
  constructor(props) {
    super(props)

    const client = new GimletCLIClient()
    client.onError = (response) => {
      console.log(response)
      console.log(`${response.status}: ${response.statusText} on ${response.path}`)
    }

    this.state = {
      client: client,
      stack: {},
      stackNonDefaultValues: {},
      toggleState: {},
    }
    this.setValues = this.setValues.bind(this)
    this.toggleComponent = this.toggleComponent.bind(this)
  }

  componentDidMount() {
    fetch('/stack-definition.json')
      .then(response => {
        if (!response.ok && window !== undefined) {
          console.log("Using fixture")
          return stackDefinitionFixture.default
        }
        return response.json()
      })
      .then(data => this.setState({stackDefinition: data}))

    // fetch('/values.schema.json')
    //   .then(response => {
    //     if (!response.ok && window !== undefined) {
    //       console.log("Using fixture")
    //       return schemaFixture.default
    //     }
    //     return response.json()
    //   })
    //   .then(data => this.setState({ schema: data }))
    //
    fetch('/stack.json')
      .then(response => {
        if (!response.ok && window !== undefined) {
          console.log("Using fixture")
          return {}
        }
        return response.json()
      })
      .then(data => this.setState({stack: data}))
  }

  setValues(variable, values, nonDefaultValues) {
    console.log(values)
    this.setState(prevState => ({
      stack: {
        ...prevState.stack,
        [variable]: values
      },
      stackNonDefaultValues: {
        ...prevState.stackNonDefaultValues,
        [variable]: nonDefaultValues
      }
    }))
    //this.state.client.saveValues(nonDefaultValues)
  }

  toggleComponent(category, component) {
    console.log("toggling " + category + " " + component)
    this.setState(prevState => ({
      toggleState: {
        ...prevState.toggleState,
        [category]: prevState.toggleState[category] == component ? undefined : component
      }
    }))
  }

  render() {
    let {stackDefinition, stack, toggleState} = this.state

    if (stackDefinition === undefined || stack === undefined) {
      return null;
    }

    const genericComponentSaver = this.setValues;
    const toggleComponentHandler = this.toggleComponent;

    return (
      <div>
        <StreamingBackend client={this.state.client}/>
        <div className="fixed bottom-0 right-0">
          <span className="inline-flex rounded-md shadow-sm m-8">
            <button
              type="button"
              className="cursor-default inline-flex items-center px-6 py-3 border border-transparent text-base leading-6 font-medium rounded-md text-white bg-gray-600 transition ease-in-out duration-150"
              onClick={() => {
                console.log(this.state.stack)
                console.log(this.state.stackNonDefaultValues)
              }}
            >
              Close the browser when you are done, the values will be printed on the console
            </button>
          </span>
        </div>
        {
          stackDefinition.categories.map(category => {
            let selectedComponent = undefined;
            let selectedComponentConfig = undefined;
            let componentSaver = undefined;
            const selectedComponentName = toggleState[category.id];
            if (selectedComponentName !== undefined) {
              const selectedComponentArray = stackDefinition.components.filter(component => component.variable === toggleState[category.id]);
              selectedComponent = selectedComponentArray[0];
              selectedComponentConfig = stack[selectedComponent.variable];
              if (selectedComponentConfig === undefined) {
                selectedComponentConfig = {}
              }
              componentSaver = function(values, nonDefaultValues) {
                genericComponentSaver(selectedComponent.variable, values, nonDefaultValues)
              };
            }

            const componentsForCategory = stackDefinition.components.filter(component => component.category === category.id);
            const componentTitles = componentsForCategory.map(component => {
              return (
                <Tile
                  category={category}
                  component={component}
                  componentConfig={stack[component.variable]}
                  selectedComponentName={selectedComponentName}
                  toggleComponentHandler={toggleComponentHandler}
                />
              )
            })

            const componentConfigPanel = selectedComponentName === undefined ? null : (
              <div>
                <div className="shadow sm:rounded-md sm:overflow-hidden">
                  <div className="bg-white py-6 px-4 space-y-6 sm:p-6">
                    <HelmUI
                      schema={selectedComponent.schema}
                      config={selectedComponent.uiSchema}
                      values={selectedComponentConfig}
                      setValues={componentSaver}
                    />
                  </div>
                </div>
              </div>
            );

            return (
              <div className="container mx-auto m-8 max-w-4xl">
                <h2>{category.name}</h2>
                <div className="flex space-x-2">
                  {componentTitles}
                </div>
                <div className="my-2">
                  {componentConfigPanel}
                </div>
              </div>
            )
          })
        }
      </div>
    )
  }
};

export default hot(module)(App)
