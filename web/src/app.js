import React, {Component} from 'react'
import {hot} from 'react-hot-loader'
import * as stackDefinitionFixture from '../fixtures/stack-definition.json'
import { XIcon } from '@heroicons/react/outline'

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

    const categories = stackDefinition.categories.map(category => {
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
        <div className="py-6 px-4 space-y-6 sm:p-6">
          <HelmUI
            schema={selectedComponent.schema}
            config={selectedComponent.uiSchema}
            values={selectedComponentConfig}
            setValues={componentSaver}
          />
        </div>
      );

      return (
        <div class="my-8">
          <h2 class="text-lg">{category.name}</h2>
          <div className="flex space-x-2 my-2">
            {componentTitles}
          </div>
          <div className="my-2">
            { selectedComponentName !== undefined &&
            <div className="px-8 py-4 shadow sm:rounded-md sm:overflow-hidden bg-white relative">
              <div className="hidden sm:block absolute top-0 right-0 pt-4 pr-4">
                <button
                  type="button"
                  className="bg-white rounded-md text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                  onClick={() => toggleComponentHandler(category.id, selectedComponent.variable)}
                >
                  <span className="sr-only">Close</span>
                  <XIcon className="h-6 w-6" aria-hidden="true" />
                </button>
              </div>
              <div>
                <div className="sm:hidden">
                  <label htmlFor="tabs" className="sr-only">Select a tab</label>
                  <select id="tabs" name="tabs"
                          className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md">
                    <option>Getting Started</option>
                    <option selected>Config</option>
                  </select>
                </div>
                <div className="hidden sm:block">
                  <div className="border-b border-gray-200">
                    <nav className="-mb-px flex space-x-8" aria-label="Tabs">
                      <a href="#"
                         className="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm">
                        Getting Started
                      </a>
                      <a href="#"
                         className="border-indigo-500 text-indigo-600 whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
                         aria-current="page"
                      >
                        Config
                      </a>
                    </nav>
                  </div>
                </div>
              </div>
              {componentConfigPanel}
            </div>
            }
          </div>
        </div>
      )
    })

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
        <div className="container mx-auto m-8 max-w-4xl">
          <h1 class="text-2xl font-bold my-16">{stackDefinition.name}
            <span class="font-normal text-lg block">{stackDefinition.description}</span>
          </h1>
          {categories}
        </div>
      </div>
    )
  }
};

export default hot(module)(App)
