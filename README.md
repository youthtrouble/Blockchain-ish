## Decentralized Weather App
This project is being built to ground my basic knowledge of crypto and the Blockchain. It is not necessarily practical 
because it is based off the assumptions that:  
- Every node that wants to  contribute, owns a thermometer.
- No node wants to replicate the [Google Maps Big Data hack](https://news.artnet.com/art-world/artist-simon-weckert-google-map-hack-1769187) in this context, in the event that the Practical Improvements
explained below are implemented.
  
### How it works locally

- Clone the repo `git clone https://github.com/youthtrouble/Blockchain-ish`
- In the terminal, run the command `go run cmd/main.go` 
  - This initialises a [genesis block](https://www.investopedia.com/terms/g/genesis-block.asp#:~:text=A%20Genesis%20Block%20is%20the,occur%20on%20a%20blockchain%20network.) on the blockchain, this blockchain can then be updated by the other nodes
on the [TCP](https://en.wikipedia.org/wiki/Transmission_Control_Protocol) network.

- Open as many terminals as you please and then connect them to the TCP network by  running the command
`nc localost 8000` The port 8000 can be changed by updating the .env file :).
  
- Fill the prompt for the Temperature Info using the format `Temperature(in degrees celsius), Location`.
    - The prompt generates a new Block based on the input and adds it to the Blockchain on the main terminal. 
    -  In 30-second Intervals, these other nodes(terminals)
   are updated with a verified copy of the Blockchain.
  

### Practical Improvements

In the event that the aforementioned assumptions are true in th real-world, ere are some ideas
 that could make this project practical:

- Sorting and analysis: The data could be sorted and analyzed over a specific time period using practical weather optimized Algorithms
 to find the median Temperature per Location.
  
- Updating an actual production Application per time period.

---

## Example

![Screenshot from 2021-06-05 17-05-29](https://user-images.githubusercontent.com/47859940/120898024-14cafd80-c621-11eb-9526-28cc4e728ff0.png)


