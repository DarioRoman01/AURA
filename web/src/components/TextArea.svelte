<script lang="ts">
import { createEventDispatcher } from "svelte";

  import { MaterialApp, Row, Col, Textarea, Button } from "svelte-materialify";
  import { parse } from "../../services/api";

  let source: string;
  const dispatch = createEventDispatcher();

  const handleSubmit = () => {
    const parsed = parse(source)
    parsed.then((evaluated) => dispatch("parse", { evaluated }))
  }
</script>

<MaterialApp>
  <Row>
    <Col>
      <Textarea solo placeholder="code..." bind:value={source}/>
    </Col>
  </Row>
  <div class="d-flex justify-center">
    <Button on:click={handleSubmit}>
      send
    </Button>
  </div>
</MaterialApp>