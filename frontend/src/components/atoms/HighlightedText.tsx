import { Filter } from "../../app/types/Filter";
import { mergeArrays, randomKey } from "../../utils/utils";
import { Text } from "./Text";

export const HighlightedText = (props: {
    text: string;
    filters?: Filter[];
}) => {
    if (!props.filters || (props.filters && props.filters.length === 0)) {
        return <Text>{props.text}</Text>
    };
    const combinedRegex = new RegExp(
        props
            .filters
            .filter(f => !f.isBlocking && f.regex)
            .map(f => `${f.regex}`).join("|"), 'g');
    const splits = props.text.split(combinedRegex);
    const matches = Array.from(props.text.matchAll(combinedRegex)).map(String);
    const result = mergeArrays(splits, matches);
    
    return (
        <div 
            className="inline-block"
        >
            {
                result.map((excerpt, idx) => {
                    return (
                        idx % 2 === 0 ? <Text key={randomKey()}>{excerpt}</Text> : (
                            <Text 
                                key={randomKey()} 
                                className="font-bold" 
                                color={props.filters?.find(f => f.regex && excerpt.match(f.regex))?.color}
                            >{excerpt}</Text>
                        )
                    )
                })
            }
        </div>
    );
}