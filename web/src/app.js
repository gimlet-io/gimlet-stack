import React, {Component} from 'react'
import {hot} from 'react-hot-loader'
import * as stackDefinitionFixture from '../fixtures/stack-definition.json'
import './style.css'
import StreamingBackend from './streamingBackend'
import GimletCLIClient from './client'

import {Category} from "./components/category";

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
    }
    this.setValues = this.setValues.bind(this)

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
    const updatedNonDefaultValues = {
      ...this.state.stackNonDefaultValues,
      [variable]: nonDefaultValues
    }

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

    this.state.client.saveValues(updatedNonDefaultValues)
  }

  render() {
    let {stackDefinition, stack} = this.state

    console.log(stackDefinition)
    console.log(stack)

    if (stackDefinition === undefined || stack === undefined) {
      return null;
    }

    const categories = stackDefinition.categories.map(category => {
      return <Category
        category={category}
        stackDefinition={stackDefinition}
        stack={stack}
        genericComponentSaver={this.setValues}
      />
    })

    console.log(categories)

    return (
      <div>
        <StreamingBackend client={this.state.client}/>
        <div className="fixed bottom-0 right-0">
          <span className="inline-flex rounded-md shadow-sm m-8">
            <button
              type="button"
              className="inline-flex items-center px-12 py-6 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              onClick={() => {
                console.log(this.state.stack)
                console.log(this.state.stackNonDefaultValues)
                close();
              }}
            >
              Close tab & <br />
              Write values
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
