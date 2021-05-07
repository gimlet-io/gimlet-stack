import './style.css'
import React, {Component} from 'react'

export class Tile extends Component {
    constructor(props) {
        super(props);
    }

    render() {
        const { category, component, componentConfig, selectedComponentName, toggleComponentHandler } = this.props;

        const selected = component.variable === selectedComponentName;
        const enabled = componentConfig !== undefined ? componentConfig.enabled : false;
        const selectedOrEnabled = selected || enabled;
        return (
            <div onClick={() => toggleComponentHandler(category.id, component.variable)}>
                <div className="w-32 h-32 px-2 overflow-hidden cursor-pointer">
                    <div className={!selectedOrEnabled ? 'bg-gray-100 hover:bg-gray-300 filter grayscale hover:grayscale-0' : 'bg-gray-300'}>
                        <img className="h-20 mx-auto pt-4" src={component.logo} alt={component.name}/>
                        <div className="font-bold text-sm py-2 text-center">{component.name}</div>
                    </div>
                </div>
            </div>
        )
    }
}
