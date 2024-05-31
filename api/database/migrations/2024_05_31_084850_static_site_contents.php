<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('static_site_contents', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->foreignId('property')->constrained(table: 'static_site_properties', column: 'property')->onUpdate('cascade')->onDelete('cascade');
            $table->foreignId('lang_code')->constrained(table: 'languages', column: 'code')->onUpdate('cascade')->onDelete('cascade');
            $table->string('content');
            $table->unique(['property','lang_code']);

            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('static_site_contents');
    }
};
