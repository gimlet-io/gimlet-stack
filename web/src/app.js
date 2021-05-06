import React, {Component} from 'react'
import {hot} from 'react-hot-loader'
import * as stackDefinitionFixture from '../fixtures/stack-definition.json'

import HelmUI from 'helm-react-ui'
import './style.css'
import StreamingBackend from './streamingBackend'
import GimletCLIClient from './client'

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
    this.setState({
      stack: {
        [variable]: values
      },
      stackNonDefaultValues: {
        [variable]: nonDefaultValues
      }
    })
    //this.state.client.saveValues(nonDefaultValues)
  }

  render() {
    let {stackDefinition, stack} = this.state

    if (stackDefinition === undefined || stack === undefined) {
      return null;
    }

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
          stackDefinition.components.map(c => {
            let componentValues = stack[c.variable];
            if (componentValues === undefined) {
              componentValues = {}
            }
            let componentSaver = function(values, nonDefaultValues) {
              this.setValues(c.variable, values, nonDefaultValues)
            }.bind(this);

            return (
              <div className="container mx-auto m-8 max-w-xl">
                <h2>{c.name}</h2>
                <div className="shadow sm:rounded-md sm:overflow-hidden">
                  <div className="bg-white py-6 px-4 space-y-6 sm:p-6">
                    <HelmUI
                      schema={c.schema}
                      config={c.uiSchema}
                      values={componentValues}
                      setValues={componentSaver}
                    />
                  </div>
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
